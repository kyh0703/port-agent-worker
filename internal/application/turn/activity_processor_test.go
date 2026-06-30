package turn

import (
	"context"
	"errors"
	"testing"
	"time"

	"port-agent-worker/internal/domain/voice"
)

func TestActivityProcessorReturnsBargeInForSpeechStartedDuringBotSpeech(t *testing.T) {
	controller := NewController(Config{StopDelay: time.Second}, nil)
	processor := NewActivityProcessor(controller)
	now := time.Unix(10, 0)
	event := mustSpeechEvent(t, voice.SpeechStarted, now)

	controller.BotStarted()
	decision, err := processor.Handle(context.Background(), event)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if !decision.BargeIn {
		t.Fatal("expected barge-in decision")
	}
	if decision.Endpoint {
		t.Fatal("expected no endpoint on speech started")
	}
}

func TestActivityProcessorChecksEndpointAfterSpeechStopped(t *testing.T) {
	controller := NewController(Config{StopDelay: time.Second}, nil)
	processor := NewActivityProcessor(controller)
	start := time.Unix(10, 0)
	event := mustSpeechEvent(t, voice.SpeechStopped, start)

	decision, err := processor.Handle(context.Background(), event)
	if err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if decision.Endpoint {
		t.Fatal("expected no endpoint at speech stopped time")
	}

	decision, err = processor.Tick(context.Background(), start.Add(time.Second))
	if err != nil {
		t.Fatalf("Tick returned error: %v", err)
	}
	if !decision.Endpoint {
		t.Fatal("expected endpoint after stop delay")
	}
}

func TestActivityProcessorRejectsInvalidEvent(t *testing.T) {
	processor := NewActivityProcessor(NewController(Config{}, nil))

	_, err := processor.Handle(context.Background(), voice.SpeechActivityEvent{})
	if !errors.Is(err, voice.ErrInvalidSpeechActivityEvent) {
		t.Fatalf("error = %v, want %v", err, voice.ErrInvalidSpeechActivityEvent)
	}
}

func TestActivityProcessorRequiresController(t *testing.T) {
	processor := NewActivityProcessor(nil)
	event := mustSpeechEvent(t, voice.SpeechStarted, time.Unix(10, 0))

	_, err := processor.Handle(context.Background(), event)
	if !errors.Is(err, ErrControllerRequired) {
		t.Fatalf("error = %v, want %v", err, ErrControllerRequired)
	}
}

func mustSpeechEvent(t *testing.T, kind voice.SpeechActivityKind, at time.Time) voice.SpeechActivityEvent {
	t.Helper()

	event, err := voice.NewSpeechActivityEvent(kind, at)
	if err != nil {
		t.Fatalf("NewSpeechActivityEvent returned error: %v", err)
	}

	return event
}
