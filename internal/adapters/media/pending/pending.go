package pending

import (
	"context"
	"errors"

	"port-agent-worker/internal/domain/voice"
)

var ErrMediaNotConfigured = errors.New("media is not configured")

type Ingress struct{}

func (Ingress) PCMFrames(context.Context) (<-chan voice.PCMFrame, error) {
	return nil, ErrMediaNotConfigured
}

type Egress struct{}

func (Egress) WritePCM(context.Context, voice.PCMFrame) error {
	return ErrMediaNotConfigured
}

func (Egress) Flush(context.Context) error {
	return ErrMediaNotConfigured
}

func (Egress) Close() error {
	return nil
}
