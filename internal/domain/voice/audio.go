package voice

import (
	"errors"
	"time"
)

type PCMEncoding string

const (
	PCMEncodingS16LE PCMEncoding = "s16le"
)

var ErrInvalidPCMFrame = errors.New("invalid pcm frame")

type PCMFrame struct {
	Data       []byte
	Encoding   PCMEncoding
	SampleRate int
	Channels   int
	Duration   time.Duration
}

func NewPCMFrame(data []byte, sampleRate int, channels int, duration time.Duration) (PCMFrame, error) {
	if len(data) == 0 || sampleRate <= 0 || channels <= 0 || duration <= 0 {
		return PCMFrame{}, ErrInvalidPCMFrame
	}

	cloned := make([]byte, len(data))
	copy(cloned, data)

	return PCMFrame{
		Data:       cloned,
		Encoding:   PCMEncodingS16LE,
		SampleRate: sampleRate,
		Channels:   channels,
		Duration:   duration,
	}, nil
}
