package pion

import (
	"errors"
	"testing"

	"github.com/pion/webrtc/v4"

	"port-agent-worker/internal/domain/voice"
)

func TestNewRuntimeRequiresTracks(t *testing.T) {
	_, err := NewRuntime(Config{Encoder: fakeEncoder{}})
	if !errors.Is(err, ErrMissingTrack) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrMissingTrack)
	}

	_, err = NewRuntime(Config{
		InputTrack: &webrtc.TrackRemote{},
		Encoder:    fakeEncoder{},
	})
	if !errors.Is(err, ErrMissingTrack) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrMissingTrack)
	}
}

func TestNewRuntimeRequiresEncoder(t *testing.T) {
	_, err := NewRuntime(Config{
		InputTrack:  &webrtc.TrackRemote{},
		OutputTrack: &webrtc.TrackLocalStaticRTP{},
	})
	if !errors.Is(err, ErrMissingEncoder) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrMissingEncoder)
	}
}

func TestNewRuntimeBuildsAudioRuntime(t *testing.T) {
	runtime, err := NewRuntime(Config{
		InputTrack:  &webrtc.TrackRemote{},
		OutputTrack: &webrtc.TrackLocalStaticRTP{},
		Encoder:     fakeEncoder{},
	})
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	if runtime.Ingress == nil {
		t.Fatal("Ingress = nil")
	}
	if runtime.Egress == nil {
		t.Fatal("Egress = nil")
	}
}

type fakeEncoder struct{}

func (fakeEncoder) Encode(voice.PCMFrame) ([]byte, error) {
	return []byte{1, 2, 3}, nil
}
