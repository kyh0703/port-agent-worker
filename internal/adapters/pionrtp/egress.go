package pionrtp

import (
	"context"

	"github.com/pion/webrtc/v4"

	"port-agent-worker/internal/domain/voice"
)

type Egress struct {
	track *webrtc.TrackLocalStaticRTP
}

func NewEgress(track *webrtc.TrackLocalStaticRTP) *Egress {
	return &Egress{track: track}
}

func (e *Egress) WritePCM(context.Context, voice.PCMFrame) error {
	if e.track == nil {
		return ErrNilTrack
	}

	return ErrMediaPipelineNotReady
}

func (e *Egress) Flush(context.Context) error {
	if e.track == nil {
		return ErrNilTrack
	}

	return nil
}

func (e *Egress) Close() error {
	return nil
}
