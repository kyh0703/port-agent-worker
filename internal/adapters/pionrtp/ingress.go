package pionrtp

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/pion/opus"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"

	"port-agent-worker/internal/domain/voice"
)

const (
	defaultSampleRate = 16000
	defaultChannels   = 1
	maxOpusSamples    = 5760
)

type Ingress struct {
	track      *webrtc.TrackRemote
	source     packetSource
	decoder    frameDecoder
	sampleRate int
	channels   int
}

func NewIngress(track *webrtc.TrackRemote) *Ingress {
	ingress := &Ingress{
		track:      track,
		sampleRate: defaultSampleRate,
		channels:   defaultChannels,
	}
	if track != nil {
		ingress.source = trackPacketSource{track: track}
		ingress.decoder = newOpusFrameDecoder(defaultSampleRate, defaultChannels)
	}
	return ingress
}

func (i *Ingress) PCMFrames(ctx context.Context) (<-chan voice.PCMFrame, error) {
	if i.track == nil {
		return nil, ErrNilTrack
	}
	if i.source == nil {
		return nil, ErrNilPacketSource
	}
	if i.decoder == nil {
		return nil, ErrNilFrameDecoder
	}

	frames := make(chan voice.PCMFrame)
	go i.readLoop(ctx, frames)
	return frames, nil
}

func (i *Ingress) readLoop(ctx context.Context, frames chan<- voice.PCMFrame) {
	defer close(frames)

	for {
		packet, err := i.source.ReadRTP(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, io.EOF) {
				return
			}
			return
		}
		if packet == nil || len(packet.Payload) == 0 {
			continue
		}

		pcm, err := i.decoder.Decode(packet.Payload)
		if err != nil || len(pcm) == 0 {
			continue
		}

		frame, err := voice.NewPCMFrame(
			int16ToBytes(pcm),
			i.sampleRate,
			i.channels,
			frameDuration(len(pcm), i.sampleRate, i.channels),
		)
		if err != nil {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case frames <- frame:
		}
	}
}

type packetSource interface {
	ReadRTP(ctx context.Context) (*rtp.Packet, error)
}

type trackPacketSource struct {
	track *webrtc.TrackRemote
}

func (s trackPacketSource) ReadRTP(ctx context.Context) (*rtp.Packet, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	packet, _, err := s.track.ReadRTP()
	if err != nil {
		return nil, err
	}
	return packet, nil
}

type frameDecoder interface {
	Decode(payload []byte) ([]int16, error)
}

type opusFrameDecoder struct {
	decoder opus.Decoder
	buffer  []int16
}

func newOpusFrameDecoder(sampleRate int, channels int) *opusFrameDecoder {
	decoder, err := opus.NewDecoderWithOutput(sampleRate, channels)
	if err != nil {
		panic(fmt.Sprintf("create opus decoder: %v", err))
	}

	return &opusFrameDecoder{
		decoder: decoder,
		buffer:  make([]int16, maxOpusSamples*channels),
	}
}

func (d *opusFrameDecoder) Decode(payload []byte) ([]int16, error) {
	samples, err := d.decoder.DecodeToInt16(payload, d.buffer)
	if err != nil {
		return nil, err
	}
	out := make([]int16, samples)
	copy(out, d.buffer[:samples])
	return out, nil
}

func int16ToBytes(samples []int16) []byte {
	out := make([]byte, len(samples)*2)
	for i, sample := range samples {
		binary.LittleEndian.PutUint16(out[i*2:], uint16(sample))
	}
	return out
}

func frameDuration(sampleCount int, sampleRate int, channels int) time.Duration {
	if sampleRate <= 0 || channels <= 0 {
		return time.Nanosecond
	}
	perChannelSamples := sampleCount / channels
	if perChannelSamples <= 0 {
		return time.Nanosecond
	}
	return time.Duration(perChannelSamples) * time.Second / time.Duration(sampleRate)
}

func newIngressFromComponents(source packetSource, decoder frameDecoder, sampleRate int, channels int) *Ingress {
	return &Ingress{
		track:      &webrtc.TrackRemote{},
		source:     source,
		decoder:    decoder,
		sampleRate: sampleRate,
		channels:   channels,
	}
}
