package session

import (
	"context"
	"errors"
	"testing"
	"time"

	"port-agent-worker/internal/application/ports"
	"port-agent-worker/internal/application/turn"
	"port-agent-worker/internal/domain/voice"
)

func TestNewTurnAwareRunnerFromRuntimeBuildsRunner(t *testing.T) {
	runner := NewTurnAwareRunnerFromRuntime(
		ProviderRuntime{
			STT: &fakeSTT{transcripts: []voice.Transcript{{Text: "hello", IsFinal: true}}},
			LLM: &fakeLLM{response: voice.AssistantResponse{Text: "hi"}},
			TTS: &fakeTTS{frames: []voice.PCMFrame{mustFrame(t, []byte{1, 2})}},
		},
		AudioRuntime{
			Ingress: &fakeIngress{frames: []voice.PCMFrame{mustFrame(t, []byte{3, 4})}},
			Egress:  &fakeEgress{},
		},
		TurnRuntime{
			VAD:       &recordingVAD{},
			Processor: turn.NewActivityProcessor(turn.NewController(turn.Config{}, nil)),
		},
	)
	if runner == nil {
		t.Fatal("runner = nil")
	}
}

func TestTurnAwareOrchestratorFansOutAudioAndHandlesTurnDecisions(t *testing.T) {
	ctx := context.Background()
	frameA := mustFrame(t, []byte{1, 2})
	frameB := mustFrame(t, []byte{3, 4})
	releaseFinal := make(chan struct{})
	decisionHandler := &recordingDecisionHandler{
		releaseAfter: 2,
		release:      releaseFinal,
	}
	stt := &recordingSTT{
		finalText:     "hello",
		releaseFinal:  releaseFinal,
		finalReleased: make(chan struct{}),
	}
	vad := &recordingVAD{
		events: []voice.SpeechActivityEvent{
			mustSpeechEvent(t, voice.SpeechStarted, time.Unix(10, 0)),
			mustSpeechEvent(t, voice.SpeechStopped, time.Unix(11, 0)),
		},
	}
	egress := &fakeEgress{}
	controller := turn.NewController(turn.Config{StopDelay: time.Nanosecond}, nil)
	controller.BotStarted()

	orchestrator := NewTurnAwareOrchestrator(
		&fakeIngress{frames: []voice.PCMFrame{frameA, frameB}},
		egress,
		stt,
		&fakeLLM{response: voice.AssistantResponse{Text: "hi"}},
		&fakeTTS{frames: []voice.PCMFrame{mustFrame(t, []byte{5, 6})}},
		TurnRuntime{
			VAD:          vad,
			Processor:    turn.NewActivityProcessor(controller),
			Handler:      decisionHandler,
			TickInterval: time.Millisecond,
		},
	)

	if err := orchestrator.RunTurn(ctx); err != nil {
		t.Fatalf("RunTurn() error = %v", err)
	}

	if stt.framesSeen != 2 {
		t.Fatalf("stt frames seen = %d, want 2", stt.framesSeen)
	}
	if vad.framesSeen != 2 {
		t.Fatalf("vad frames seen = %d, want 2", vad.framesSeen)
	}
	if !decisionHandler.seenBargeIn {
		t.Fatal("expected barge-in decision")
	}
	if !decisionHandler.seenEndpoint {
		t.Fatal("expected endpoint decision")
	}
	if !egress.flushed {
		t.Fatal("expected egress flush")
	}
}

func TestTurnAwareOrchestratorReturnsTurnDecisionHandlerError(t *testing.T) {
	handlerErr := errors.New("handler failed")
	controller := turn.NewController(turn.Config{}, nil)
	controller.BotStarted()
	orchestrator := NewTurnAwareOrchestrator(
		&fakeIngress{frames: []voice.PCMFrame{mustFrame(t, []byte{1, 2})}},
		&fakeEgress{},
		&recordingSTT{releaseFinal: make(chan struct{})},
		&fakeLLM{},
		&fakeTTS{},
		TurnRuntime{
			VAD: &recordingVAD{events: []voice.SpeechActivityEvent{
				mustSpeechEvent(t, voice.SpeechStarted, time.Unix(10, 0)),
			}},
			Processor: turn.NewActivityProcessor(controller),
			Handler:   failingDecisionHandler{err: handlerErr},
		},
	)

	err := orchestrator.RunTurn(context.Background())
	if !errors.Is(err, handlerErr) {
		t.Fatalf("RunTurn() error = %v, want %v", err, handlerErr)
	}
}

func TestTurnAwareOrchestratorReturnsVADError(t *testing.T) {
	vadErr := errors.New("vad failed")
	orchestrator := NewTurnAwareOrchestrator(
		&fakeIngress{frames: []voice.PCMFrame{mustFrame(t, []byte{1, 2})}},
		&fakeEgress{},
		&fakeSTT{transcripts: []voice.Transcript{{Text: "hello", IsFinal: true}}},
		&fakeLLM{},
		&fakeTTS{},
		TurnRuntime{
			VAD:       failingVAD{err: vadErr},
			Processor: turn.NewActivityProcessor(turn.NewController(turn.Config{}, nil)),
		},
	)

	err := orchestrator.RunTurn(context.Background())
	if !errors.Is(err, vadErr) {
		t.Fatalf("RunTurn() error = %v, want %v", err, vadErr)
	}
}

func TestTurnAwareOrchestratorCompletesWhenVADStopsReading(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	frames := []voice.PCMFrame{
		mustFrame(t, []byte{1, 2}),
		mustFrame(t, []byte{3, 4}),
		mustFrame(t, []byte{5, 6}),
	}
	egress := &fakeEgress{}
	stt := &recordingSTT{
		finalText:     "hello",
		releaseFinal:  closedSignal(),
		finalReleased: make(chan struct{}),
	}

	orchestrator := NewTurnAwareOrchestrator(
		&fakeIngress{frames: frames},
		egress,
		stt,
		&fakeLLM{response: voice.AssistantResponse{Text: "hi"}},
		&fakeTTS{frames: []voice.PCMFrame{mustFrame(t, []byte{7, 8})}},
		TurnRuntime{
			VAD:          earlyClosingVAD{},
			Processor:    turn.NewActivityProcessor(turn.NewController(turn.Config{}, nil)),
			TickInterval: time.Millisecond,
		},
	)

	if err := orchestrator.RunTurn(ctx); err != nil {
		t.Fatalf("RunTurn() error = %v", err)
	}
	if stt.framesSeen != len(frames) {
		t.Fatalf("stt frames seen = %d, want %d", stt.framesSeen, len(frames))
	}
	if !egress.flushed {
		t.Fatal("expected egress flush")
	}
}

var _ ports.SpeechToText = (*recordingSTT)(nil)
var _ ports.VoiceActivityDetector = (*recordingVAD)(nil)
var _ TurnDecisionHandler = (*recordingDecisionHandler)(nil)

type recordingSTT struct {
	finalText     string
	releaseFinal  <-chan struct{}
	finalReleased chan struct{}
	framesSeen    int
}

func (f *recordingSTT) Transcribe(ctx context.Context, frames <-chan voice.PCMFrame) (<-chan voice.Transcript, error) {
	out := make(chan voice.Transcript, 1)
	go func() {
		defer close(out)
		for range frames {
			f.framesSeen++
		}

		select {
		case <-ctx.Done():
			return
		case <-f.releaseFinal:
		}

		out <- voice.Transcript{Text: f.finalText, IsFinal: true}
		close(f.finalReleased)
	}()
	return out, nil
}

type recordingVAD struct {
	events     []voice.SpeechActivityEvent
	framesSeen int
}

type failingVAD struct {
	err error
}

type earlyClosingVAD struct{}

func (f failingVAD) DetectSpeech(context.Context, <-chan voice.PCMFrame) (<-chan voice.SpeechActivityEvent, error) {
	return nil, f.err
}

func (earlyClosingVAD) DetectSpeech(context.Context, <-chan voice.PCMFrame) (<-chan voice.SpeechActivityEvent, error) {
	out := make(chan voice.SpeechActivityEvent)
	close(out)
	return out, nil
}

func (f *recordingVAD) DetectSpeech(ctx context.Context, frames <-chan voice.PCMFrame) (<-chan voice.SpeechActivityEvent, error) {
	out := make(chan voice.SpeechActivityEvent, len(f.events))
	go func() {
		defer close(out)
		for range frames {
			f.framesSeen++
		}

		for _, event := range f.events {
			select {
			case <-ctx.Done():
				return
			case out <- event:
			}
		}
	}()
	return out, nil
}

type recordingDecisionHandler struct {
	releaseAfter int
	release      chan<- struct{}
	released     bool
	seenBargeIn  bool
	seenEndpoint bool
}

type failingDecisionHandler struct {
	err error
}

func (h failingDecisionHandler) HandleTurnDecision(context.Context, turn.Decision) error {
	return h.err
}

func (h *recordingDecisionHandler) HandleTurnDecision(_ context.Context, decision turn.Decision) error {
	if decision.BargeIn {
		h.seenBargeIn = true
	}
	if decision.Endpoint {
		h.seenEndpoint = true
	}

	if !h.released && h.releaseAfter <= h.count() {
		close(h.release)
		h.released = true
	}
	return nil
}

func (h *recordingDecisionHandler) count() int {
	count := 0
	if h.seenBargeIn {
		count++
	}
	if h.seenEndpoint {
		count++
	}
	return count
}

func closedSignal() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func mustSpeechEvent(t *testing.T, kind voice.SpeechActivityKind, at time.Time) voice.SpeechActivityEvent {
	t.Helper()

	event, err := voice.NewSpeechActivityEvent(kind, at)
	if err != nil {
		t.Fatalf("NewSpeechActivityEvent() error = %v", err)
	}
	return event
}
