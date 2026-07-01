package pion

import (
	"errors"
	"fmt"

	"github.com/pion/webrtc/v4"

	"port-voice-pipeline/internal/adapters/pionrtp"
	"port-voice-pipeline/internal/application/session"
)

var ErrMissingTrack = errors.New("pion media track is required")
var ErrMissingEncoder = errors.New("pion egress encoder is required")

type Config struct {
	InputTrack  *webrtc.TrackRemote
	OutputTrack *webrtc.TrackLocalStaticRTP
	Encoder     pionrtp.FrameEncoder
}

func NewRuntime(config Config) (session.AudioRuntime, error) {
	if config.InputTrack == nil {
		return session.AudioRuntime{}, fmt.Errorf("%w: input", ErrMissingTrack)
	}
	if config.OutputTrack == nil {
		return session.AudioRuntime{}, fmt.Errorf("%w: output", ErrMissingTrack)
	}
	if config.Encoder == nil {
		return session.AudioRuntime{}, ErrMissingEncoder
	}

	return session.AudioRuntime{
		Ingress: pionrtp.NewIngress(config.InputTrack),
		Egress:  pionrtp.NewEgressWithEncoder(config.OutputTrack, config.Encoder),
	}, nil
}
