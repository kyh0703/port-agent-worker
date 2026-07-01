package noop

import (
	"context"
	"testing"
	"time"

	"port-voice-pipeline/internal/domain/voice"
)

func TestDetectorDrainsFramesAndClosesEvents(t *testing.T) {
	frame := mustFrame(t)
	frames := make(chan voice.PCMFrame, 2)
	frames <- frame
	frames <- frame
	close(frames)

	events, err := Detector{}.DetectSpeech(context.Background(), frames)
	if err != nil {
		t.Fatalf("DetectSpeech() error = %v", err)
	}

	if _, ok := <-events; ok {
		t.Fatal("expected no speech activity events")
	}
}

func TestDetectorStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	frames := make(chan voice.PCMFrame)
	events, err := Detector{}.DetectSpeech(ctx, frames)
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
		t.Fatal("timed out waiting for events channel close")
	}
}

func mustFrame(t *testing.T) voice.PCMFrame {
	t.Helper()

	frame, err := voice.NewPCMFrame([]byte{1, 2}, 16000, 1, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("NewPCMFrame() error = %v", err)
	}
	return frame
}
