package pending

import (
	"context"
	"errors"
	"testing"
	"time"

	"port-agent-worker/internal/domain/voice"
)

func TestIngressReturnsMediaNotConfigured(t *testing.T) {
	_, err := Ingress{}.PCMFrames(context.Background())
	if !errors.Is(err, ErrMediaNotConfigured) {
		t.Fatalf("PCMFrames() error = %v, want %v", err, ErrMediaNotConfigured)
	}
}

func TestEgressReturnsMediaNotConfigured(t *testing.T) {
	frame, err := voice.NewPCMFrame([]byte{1, 2}, 16000, 1, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("NewPCMFrame() error = %v", err)
	}

	err = Egress{}.WritePCM(context.Background(), frame)
	if !errors.Is(err, ErrMediaNotConfigured) {
		t.Fatalf("WritePCM() error = %v, want %v", err, ErrMediaNotConfigured)
	}

	err = Egress{}.Flush(context.Background())
	if !errors.Is(err, ErrMediaNotConfigured) {
		t.Fatalf("Flush() error = %v, want %v", err, ErrMediaNotConfigured)
	}
}
