package pionrtp

import (
	"context"

	"github.com/pion/webrtc/v4"

	"port-agent-worker/internal/domain/voice"
)

type Ingress struct {
	track *webrtc.TrackRemote
}

func NewIngress(track *webrtc.TrackRemote) *Ingress {
	return &Ingress{track: track}
}

func (i *Ingress) PCMFrames(context.Context) (<-chan voice.PCMFrame, error) {
	if i.track == nil {
		return nil, ErrNilTrack
	}

	return nil, ErrMediaPipelineNotReady
}
