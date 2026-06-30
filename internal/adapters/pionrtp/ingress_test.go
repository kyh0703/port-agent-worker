package pionrtp

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/pion/rtp"
)

func TestIngressRequiresTrack(t *testing.T) {
	_, err := NewIngress(nil).PCMFrames(context.Background())
	if !errors.Is(err, ErrNilTrack) {
		t.Fatalf("PCMFrames() error = %v, want %v", err, ErrNilTrack)
	}
}

func TestIngressConvertsRTPPayloadToPCMFrame(t *testing.T) {
	source := &fakePacketSource{
		packets: []*rtp.Packet{
			{Payload: []byte{1, 2, 3}},
		},
	}
	decoder := &fakeFrameDecoder{
		samples: []int16{1, -2, 300},
	}
	ingress := newIngressFromComponents(source, decoder, 16000, 1)

	frames, err := ingress.PCMFrames(context.Background())
	if err != nil {
		t.Fatalf("PCMFrames() error = %v", err)
	}

	frame, ok := <-frames
	if !ok {
		t.Fatal("frames channel closed before first frame")
	}

	want := []byte{1, 0, 254, 255, 44, 1}
	if string(frame.Data) != string(want) {
		t.Fatalf("frame data = %v, want %v", frame.Data, want)
	}
	if frame.SampleRate != 16000 || frame.Channels != 1 {
		t.Fatalf("frame metadata = %+v", frame)
	}
	if frame.Duration != 187500*time.Nanosecond {
		t.Fatalf("frame duration = %s, want 187.5us", frame.Duration)
	}
}

func TestIngressStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ingress := newIngressFromComponents(&blockingPacketSource{}, &fakeFrameDecoder{}, 16000, 1)

	frames, err := ingress.PCMFrames(ctx)
	if err != nil {
		t.Fatalf("PCMFrames() error = %v", err)
	}

	select {
	case _, ok := <-frames:
		if ok {
			t.Fatal("frames channel open after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for frames channel close")
	}
}

type fakePacketSource struct {
	packets []*rtp.Packet
	index   int
}

func (s *fakePacketSource) ReadRTP(context.Context) (*rtp.Packet, error) {
	if s.index >= len(s.packets) {
		return nil, io.EOF
	}
	packet := s.packets[s.index]
	s.index++
	return packet, nil
}

type blockingPacketSource struct{}

func (blockingPacketSource) ReadRTP(ctx context.Context) (*rtp.Packet, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

type fakeFrameDecoder struct {
	samples []int16
}

func (d *fakeFrameDecoder) Decode([]byte) ([]int16, error) {
	return d.samples, nil
}
