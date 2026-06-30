package pionrtp

import "errors"

var (
	ErrNilTrack              = errors.New("nil pion track")
	ErrMediaPipelineNotReady = errors.New("pion rtp media pipeline is not implemented")
)
