package session

import (
	"context"
	"errors"
	"fmt"

	"port-agent-worker/internal/application/ports"
	"port-agent-worker/internal/domain/voice"
)

var ErrNoFinalTranscript = errors.New("no final transcript")

type Orchestrator struct {
	ingress ports.AudioIngress
	egress  ports.AudioEgress
	stt     ports.SpeechToText
	llm     ports.LanguageModel
	tts     ports.TextToSpeech
}

func NewOrchestrator(
	ingress ports.AudioIngress,
	egress ports.AudioEgress,
	stt ports.SpeechToText,
	llm ports.LanguageModel,
	tts ports.TextToSpeech,
) *Orchestrator {
	return &Orchestrator{
		ingress: ingress,
		egress:  egress,
		stt:     stt,
		llm:     llm,
		tts:     tts,
	}
}

func (o *Orchestrator) RunTurn(ctx context.Context) error {
	audio, err := o.ingress.PCMFrames(ctx)
	if err != nil {
		return fmt.Errorf("open audio ingress: %w", err)
	}

	transcripts, err := o.stt.Transcribe(ctx, audio)
	if err != nil {
		return fmt.Errorf("start stt: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case transcript, ok := <-transcripts:
			if !ok {
				return ErrNoFinalTranscript
			}

			text, ok := transcript.FinalText()
			if !ok {
				continue
			}

			return o.respond(ctx, voice.UserUtterance{Text: text})
		}
	}
}

func (o *Orchestrator) respond(ctx context.Context, utterance voice.UserUtterance) error {
	response, err := o.llm.Generate(ctx, utterance)
	if err != nil {
		return fmt.Errorf("generate response: %w", err)
	}

	audio, err := o.tts.Synthesize(ctx, response)
	if err != nil {
		return fmt.Errorf("start tts: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case frame, ok := <-audio:
			if !ok {
				if err := o.egress.Flush(ctx); err != nil {
					return fmt.Errorf("flush audio egress: %w", err)
				}
				return nil
			}

			if err := o.egress.WritePCM(ctx, frame); err != nil {
				return fmt.Errorf("write audio egress: %w", err)
			}
		}
	}
}
