package session

import (
	"context"
	"fmt"
	"time"

	"port-agent-worker/internal/application/ports"
	"port-agent-worker/internal/application/turn"
	"port-agent-worker/internal/domain/voice"
)

const defaultTurnTickInterval = 50 * time.Millisecond

type TurnDecisionHandler interface {
	HandleTurnDecision(ctx context.Context, decision turn.Decision) error
}

type TurnRuntime struct {
	VAD          ports.VoiceActivityDetector
	Processor    *turn.ActivityProcessor
	Handler      TurnDecisionHandler
	TickInterval time.Duration
}

type TurnAwareOrchestrator struct {
	ingress ports.AudioIngress
	egress  ports.AudioEgress
	stt     ports.SpeechToText
	llm     ports.LanguageModel
	tts     ports.TextToSpeech
	turn    TurnRuntime
}

func NewTurnAwareOrchestrator(
	ingress ports.AudioIngress,
	egress ports.AudioEgress,
	stt ports.SpeechToText,
	llm ports.LanguageModel,
	tts ports.TextToSpeech,
	turnRuntime TurnRuntime,
) *TurnAwareOrchestrator {
	if turnRuntime.TickInterval <= 0 {
		turnRuntime.TickInterval = defaultTurnTickInterval
	}

	return &TurnAwareOrchestrator{
		ingress: ingress,
		egress:  egress,
		stt:     stt,
		llm:     llm,
		tts:     tts,
		turn:    turnRuntime,
	}
}

func (o *TurnAwareOrchestrator) RunTurn(ctx context.Context) error {
	if o.turn.VAD == nil || o.turn.Processor == nil {
		return NewOrchestrator(o.ingress, o.egress, o.stt, o.llm, o.tts).RunTurn(ctx)
	}

	audio, err := o.ingress.PCMFrames(ctx)
	if err != nil {
		return fmt.Errorf("open audio ingress: %w", err)
	}

	turnCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	sttAudio := make(chan voice.PCMFrame)
	vadAudio := make(chan voice.PCMFrame)
	go fanOutPCM(turnCtx, audio, sttAudio, vadAudio)

	transcripts, err := o.stt.Transcribe(turnCtx, sttAudio)
	if err != nil {
		return fmt.Errorf("start stt: %w", err)
	}

	activity, err := o.turn.VAD.DetectSpeech(turnCtx, vadAudio)
	if err != nil {
		return fmt.Errorf("start vad: %w", err)
	}

	ticker := time.NewTicker(o.turn.TickInterval)
	defer ticker.Stop()
	tickC := ticker.C

	endpointNotified := false
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-activity:
			if !ok {
				activity = nil
				continue
			}
			decision, err := o.turn.Processor.Handle(turnCtx, event)
			if err != nil {
				return fmt.Errorf("process activity event: %w", err)
			}
			notified, err := o.handleDecision(turnCtx, decision, endpointNotified)
			if err != nil {
				return err
			}
			endpointNotified = endpointNotified || notified
			if !endpointNotified {
				decision, err := o.turn.Processor.Tick(turnCtx, time.Now())
				if err != nil {
					return fmt.Errorf("tick turn processor: %w", err)
				}
				notified, err := o.handleDecision(turnCtx, decision, endpointNotified)
				if err != nil {
					return err
				}
				endpointNotified = endpointNotified || notified
			}
			if endpointNotified {
				tickC = nil
			}
		case <-tickC:
			decision, err := o.turn.Processor.Tick(turnCtx, time.Now())
			if err != nil {
				return fmt.Errorf("tick turn processor: %w", err)
			}
			notified, err := o.handleDecision(turnCtx, decision, endpointNotified)
			if err != nil {
				return err
			}
			endpointNotified = endpointNotified || notified
			if endpointNotified {
				tickC = nil
			}
		case transcript, ok := <-transcripts:
			if !ok {
				return ErrNoFinalTranscript
			}

			text, ok := transcript.FinalText()
			if !ok {
				continue
			}

			return o.respond(turnCtx, voice.UserUtterance{Text: text})
		}
	}
}

func (o *TurnAwareOrchestrator) respond(ctx context.Context, utterance voice.UserUtterance) error {
	return NewOrchestrator(o.ingress, o.egress, o.stt, o.llm, o.tts).respond(ctx, utterance)
}

func (o *TurnAwareOrchestrator) handleDecision(ctx context.Context, decision turn.Decision, endpointNotified bool) (bool, error) {
	if !decision.BargeIn && !decision.Endpoint {
		return false, nil
	}
	if decision.Endpoint && endpointNotified {
		return false, nil
	}
	if o.turn.Handler == nil {
		return decision.Endpoint, nil
	}
	if err := o.turn.Handler.HandleTurnDecision(ctx, decision); err != nil {
		return false, fmt.Errorf("handle turn decision: %w", err)
	}
	return decision.Endpoint, nil
}

func fanOutPCM(ctx context.Context, input <-chan voice.PCMFrame, outputs ...chan<- voice.PCMFrame) {
	defer func() {
		for _, output := range outputs {
			close(output)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case frame, ok := <-input:
			if !ok {
				return
			}
			for _, output := range outputs {
				select {
				case <-ctx.Done():
					return
				case output <- frame:
				}
			}
		}
	}
}
