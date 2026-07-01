package ports

import (
	"context"

	"port-voice-pipeline/internal/domain/voice"
)

type AudioIngress interface {
	PCMFrames(ctx context.Context) (<-chan voice.PCMFrame, error)
}

type AudioEgress interface {
	WritePCM(ctx context.Context, frame voice.PCMFrame) error
	Flush(ctx context.Context) error
	Close() error
}
