package pion

import (
	"context"
	"errors"
	"testing"

	"github.com/pion/webrtc/v4"

	"port-agent-worker/internal/application/session"
	"port-agent-worker/internal/domain/voice"
)

func TestNewRunnerRequiresMediaRuntime(t *testing.T) {
	_, err := NewRunner(Config{}, session.ProviderRuntime{})
	if !errors.Is(err, ErrMissingTrack) {
		t.Fatalf("NewRunner() error = %v, want %v", err, ErrMissingTrack)
	}
}

func TestNewRunnerBuildsRunner(t *testing.T) {
	runner, err := NewRunner(
		Config{
			InputTrack:  &webrtc.TrackRemote{},
			OutputTrack: &webrtc.TrackLocalStaticRTP{},
			Encoder:     fakeEncoder{},
		},
		session.ProviderRuntime{
			STT: fakeSTT{},
			LLM: fakeLLM{},
			TTS: fakeTTS{},
		},
	)
	if err != nil {
		t.Fatalf("NewRunner() error = %v", err)
	}
	if runner == nil {
		t.Fatal("runner = nil")
	}
}

type fakeSTT struct{}

func (fakeSTT) Transcribe(context.Context, <-chan voice.PCMFrame) (<-chan voice.Transcript, error) {
	return make(chan voice.Transcript), nil
}

type fakeLLM struct{}

func (fakeLLM) Generate(context.Context, voice.UserUtterance) (voice.AssistantResponse, error) {
	return voice.AssistantResponse{Text: "ok"}, nil
}

type fakeTTS struct{}

func (fakeTTS) Synthesize(context.Context, voice.AssistantResponse) (<-chan voice.PCMFrame, error) {
	return make(chan voice.PCMFrame), nil
}
