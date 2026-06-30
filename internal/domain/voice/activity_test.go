package voice

import (
	"errors"
	"testing"
	"time"
)

func TestNewSpeechActivityEventAcceptsValidKinds(t *testing.T) {
	at := time.Unix(10, 0)

	tests := []struct {
		name string
		kind SpeechActivityKind
	}{
		{name: "speech started", kind: SpeechStarted},
		{name: "speech stopped", kind: SpeechStopped},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := NewSpeechActivityEvent(tt.kind, at)
			if err != nil {
				t.Fatalf("NewSpeechActivityEvent returned error: %v", err)
			}
			if event.Kind != tt.kind {
				t.Fatalf("kind = %v, want %v", event.Kind, tt.kind)
			}
			if !event.At.Equal(at) {
				t.Fatalf("at = %s, want %s", event.At, at)
			}
		})
	}
}

func TestNewSpeechActivityEventRejectsInvalidKind(t *testing.T) {
	_, err := NewSpeechActivityEvent(SpeechActivityKind(0), time.Unix(10, 0))
	if !errors.Is(err, ErrInvalidSpeechActivityEvent) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidSpeechActivityEvent)
	}
}

func TestNewSpeechActivityEventRejectsZeroTime(t *testing.T) {
	_, err := NewSpeechActivityEvent(SpeechStarted, time.Time{})
	if !errors.Is(err, ErrInvalidSpeechActivityEvent) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidSpeechActivityEvent)
	}
}
