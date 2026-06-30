package pionrtp

import (
	"context"
	"sync"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"

	"port-agent-worker/internal/domain/voice"
)

const (
	defaultOpusPayloadType = 111
	defaultRTPClockRate    = 48000
	defaultSSRC            = 1
)

type Egress struct {
	track       *webrtc.TrackLocalStaticRTP
	writer      rtpWriter
	encoder     FrameEncoder
	payloadType uint8
	ssrc        uint32
	clockRate   uint32
	mu          sync.Mutex
	sequence    uint16
	timestamp   uint32
}

func NewEgress(track *webrtc.TrackLocalStaticRTP) *Egress {
	return NewEgressWithEncoder(track, nil)
}

func NewEgressWithEncoder(track *webrtc.TrackLocalStaticRTP, encoder FrameEncoder) *Egress {
	egress := newEgress(trackWriter{track: track}, encoder, EgressConfig{})
	egress.track = track
	return egress
}

type EgressConfig struct {
	PayloadType uint8
	SSRC        uint32
	ClockRate   uint32
}

func newEgress(writer rtpWriter, encoder FrameEncoder, config EgressConfig) *Egress {
	if config.PayloadType == 0 {
		config.PayloadType = defaultOpusPayloadType
	}
	if config.SSRC == 0 {
		config.SSRC = defaultSSRC
	}
	if config.ClockRate == 0 {
		config.ClockRate = defaultRTPClockRate
	}

	return &Egress{
		writer:      writer,
		encoder:     encoder,
		payloadType: config.PayloadType,
		ssrc:        config.SSRC,
		clockRate:   config.ClockRate,
	}
}

func (e *Egress) WritePCM(ctx context.Context, frame voice.PCMFrame) error {
	if e.writer == nil {
		return ErrNilRTPWriter
	}
	if e.encoder == nil {
		return ErrMediaPipelineNotReady
	}

	payload, err := e.encoder.Encode(frame)
	if err != nil {
		return err
	}

	packet := e.nextPacket(payload, frame.Duration)
	return e.writer.WriteRTP(ctx, packet)
}

func (e *Egress) nextPacket(payload []byte, duration time.Duration) *rtp.Packet {
	e.mu.Lock()
	defer e.mu.Unlock()

	packet := &rtp.Packet{
		Header: rtp.Header{
			Version:        2,
			PayloadType:    e.payloadType,
			SequenceNumber: e.sequence,
			Timestamp:      e.timestamp,
			SSRC:           e.ssrc,
		},
		Payload: payload,
	}

	e.sequence++
	e.timestamp += rtpTimestampIncrement(duration, e.clockRate)
	return packet
}

func rtpTimestampIncrement(duration time.Duration, clockRate uint32) uint32 {
	if duration <= 0 || clockRate == 0 {
		return 0
	}
	return uint32(duration.Nanoseconds() * int64(clockRate) / int64(time.Second))
}

type rtpWriter interface {
	WriteRTP(ctx context.Context, packet *rtp.Packet) error
}

type trackWriter struct {
	track *webrtc.TrackLocalStaticRTP
}

func (w trackWriter) WriteRTP(ctx context.Context, packet *rtp.Packet) error {
	if w.track == nil {
		return ErrNilTrack
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return w.track.WriteRTP(packet)
}

type FrameEncoder interface {
	Encode(frame voice.PCMFrame) ([]byte, error)
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
