package pion

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pion/webrtc/v4"

	"port-agent-worker/internal/application/session"
	"port-agent-worker/internal/application/turn"
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

func TestNewTurnAwareRunnerRequiresMediaRuntime(t *testing.T) {
	_, err := NewTurnAwareRunner(Config{}, session.ProviderRuntime{}, session.TurnRuntime{})
	if !errors.Is(err, ErrMissingTrack) {
		t.Fatalf("NewTurnAwareRunner() error = %v, want %v", err, ErrMissingTrack)
	}
}

func TestNewTurnAwareRunnerBuildsRunner(t *testing.T) {
	runner, err := NewTurnAwareRunner(
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
		session.TurnRuntime{
			VAD:          fakeVAD{},
			Processor:    turn.NewActivityProcessor(turn.NewController(turn.Config{StopDelay: time.Second}, nil)),
			TickInterval: time.Second,
		},
	)
	if err != nil {
		t.Fatalf("NewTurnAwareRunner() error = %v", err)
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

type fakeVAD struct{}

func (fakeVAD) DetectSpeech(context.Context, <-chan voice.PCMFrame) (<-chan voice.SpeechActivityEvent, error) {
	events := make(chan voice.SpeechActivityEvent)
	close(events)
	return events, nil
}
