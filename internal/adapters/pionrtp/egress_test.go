package pionrtp

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pion/rtp"

	"port-agent-worker/internal/domain/voice"
)

func TestEgressRequiresTrack(t *testing.T) {
	frame := mustPCMFrame(t)

	err := NewEgressWithEncoder(nil, fakeFrameEncoder{payload: []byte{1}}).WritePCM(context.Background(), frame)
	if !errors.Is(err, ErrNilTrack) {
		t.Fatalf("WritePCM() error = %v, want %v", err, ErrNilTrack)
	}
}

func TestEgressRequiresEncoder(t *testing.T) {
	frame := mustPCMFrame(t)
	egress := newEgress(&fakeRTPWriter{}, nil, EgressConfig{})

	err := egress.WritePCM(context.Background(), frame)
	if !errors.Is(err, ErrMediaPipelineNotReady) {
		t.Fatalf("WritePCM() error = %v, want %v", err, ErrMediaPipelineNotReady)
	}
}

func TestEgressPacketizesEncodedPayload(t *testing.T) {
	writer := &fakeRTPWriter{}
	encoder := fakeFrameEncoder{payload: []byte{9, 8, 7}}
	egress := newEgress(writer, encoder, EgressConfig{
		PayloadType: 111,
		SSRC:        42,
		ClockRate:   48000,
	})

	if err := egress.WritePCM(context.Background(), mustPCMFrame(t)); err != nil {
		t.Fatalf("WritePCM() error = %v", err)
	}
	if err := egress.WritePCM(context.Background(), mustPCMFrame(t)); err != nil {
		t.Fatalf("second WritePCM() error = %v", err)
	}

	if len(writer.packets) != 2 {
		t.Fatalf("packet count = %d, want 2", len(writer.packets))
	}

	first := writer.packets[0]
	if first.Version != 2 || first.PayloadType != 111 || first.SequenceNumber != 0 || first.Timestamp != 0 || first.SSRC != 42 {
		t.Fatalf("first RTP header = %+v", first.Header)
	}
	if string(first.Payload) != string([]byte{9, 8, 7}) {
		t.Fatalf("first payload = %v", first.Payload)
	}

	second := writer.packets[1]
	if second.SequenceNumber != 1 {
		t.Fatalf("second sequence = %d, want 1", second.SequenceNumber)
	}
	if second.Timestamp != 960 {
		t.Fatalf("second timestamp = %d, want 960", second.Timestamp)
	}
}

func TestRTPTimestampIncrement(t *testing.T) {
	got := rtpTimestampIncrement(20*time.Millisecond, 48000)
	if got != 960 {
		t.Fatalf("rtpTimestampIncrement() = %d, want 960", got)
	}
}

func mustPCMFrame(t *testing.T) voice.PCMFrame {
	t.Helper()

	frame, err := voice.NewPCMFrame([]byte{1, 2, 3, 4}, 16000, 1, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("NewPCMFrame() error = %v", err)
	}
	return frame
}

type fakeRTPWriter struct {
	packets []*rtp.Packet
}

func (w *fakeRTPWriter) WriteRTP(_ context.Context, packet *rtp.Packet) error {
	w.packets = append(w.packets, packet)
	return nil
}

type fakeFrameEncoder struct {
	payload []byte
}

func (e fakeFrameEncoder) Encode(voice.PCMFrame) ([]byte, error) {
	return e.payload, nil
}
