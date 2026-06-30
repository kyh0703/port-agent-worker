package pionrtp

import "errors"

var (
	ErrNilTrack              = errors.New("nil pion track")
	ErrNilPacketSource       = errors.New("nil rtp packet source")
	ErrNilFrameDecoder       = errors.New("nil audio frame decoder")
	ErrNilRTPWriter          = errors.New("nil rtp writer")
	ErrMediaPipelineNotReady = errors.New("pion rtp media pipeline is not implemented")
)
