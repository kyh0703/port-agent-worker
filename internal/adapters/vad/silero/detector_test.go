package silero

import (
	"context"
	"errors"
	"testing"
	"time"

	"port-agent-worker/internal/application/ports"
	"port-agent-worker/internal/domain/voice"
)

func TestNewRejectsMissingEngine(t *testing.T) {
	_, err := New(nil, Config{})
	if !errors.Is(err, ErrEngineRequired) {
		t.Fatalf("New() error = %v, want %v", err, ErrEngineRequired)
	}
}

func TestNewRejectsInvalidConfig(t *testing.T) {
	_, err := New(&sequenceEngine{}, Config{Threshold: 1.5})
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("New() error = %v, want %v", err, ErrInvalidConfig)
	}
}

func TestDetectorEmitsSpeechStartedAndStopped(t *testing.T) {
	detector, err := New(&sequenceEngine{
		probabilities: []float64{0.1, 0.9, 0.8, 0.1, 0.0},
	}, Config{
		Threshold:        0.5,
		MinSpeechFrames:  2,
		MinSilenceFrames: 2,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	events := detect(t, detector, 5)

	if len(events) != 2 {
		t.Fatalf("events = %d, want 2", len(events))
	}
	if events[0].Kind != voice.SpeechStarted {
		t.Fatalf("first event = %v, want SpeechStarted", events[0].Kind)
	}
	if events[1].Kind != voice.SpeechStopped {
		t.Fatalf("second event = %v, want SpeechStopped", events[1].Kind)
	}
}

func TestDetectorDoesNotEmitDuplicateSpeechStarted(t *testing.T) {
	detector, err := New(&sequenceEngine{
		probabilities: []float64{0.9, 0.9, 0.8, 0.7},
	}, Config{
		Threshold:        0.5,
		MinSpeechFrames:  1,
		MinSilenceFrames: 1,
	})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	events := detect(t, detector, 4)

	if len(events) != 1 {
		t.Fatalf("events = %d, want 1", len(events))
	}
	if events[0].Kind != voice.SpeechStarted {
		t.Fatalf("event = %v, want SpeechStarted", events[0].Kind)
	}
}

func TestDetectorClosesEventsOnContextCancel(t *testing.T) {
	detector, err := New(blockingEngine{}, Config{})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	frames := make(chan voice.PCMFrame)
	events, err := detector.DetectSpeech(ctx, frames)
	if err != nil {
		t.Fatalf("DetectSpeech() error = %v", err)
	}

	cancel()

	select {
	case _, ok := <-events:
		if ok {
			t.Fatal("expected closed events channel")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event channel close")
	}
}

var _ ports.VoiceActivityDetector = (*Detector)(nil)

func detect(t *testing.T, detector *Detector, frameCount int) []voice.SpeechActivityEvent {
	t.Helper()

	frames := make(chan voice.PCMFrame, frameCount)
	for i := 0; i < frameCount; i++ {
		frames <- mustFrame(t)
	}
	close(frames)

	events, err := detector.DetectSpeech(context.Background(), frames)
	if err != nil {
		t.Fatalf("DetectSpeech() error = %v", err)
	}

	var collected []voice.SpeechActivityEvent
	for event := range events {
		collected = append(collected, event)
	}
	return collected
}

func mustFrame(t *testing.T) voice.PCMFrame {
	t.Helper()

	frame, err := voice.NewPCMFrame([]byte{1, 2}, 16000, 1, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("NewPCMFrame() error = %v", err)
	}
	return frame
}

type sequenceEngine struct {
	probabilities []float64
	index         int
}

func (e *sequenceEngine) SpeechProbability(context.Context, voice.PCMFrame) (float64, error) {
	if e.index >= len(e.probabilities) {
		return 0, nil
	}
	probability := e.probabilities[e.index]
	e.index++
	return probability, nil
}

type blockingEngine struct{}

func (blockingEngine) SpeechProbability(ctx context.Context, _ voice.PCMFrame) (float64, error) {
	<-ctx.Done()
	return 0, ctx.Err()
}
