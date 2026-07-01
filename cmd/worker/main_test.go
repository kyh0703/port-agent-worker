package main

import (
	"context"
	"errors"
	"testing"

	"port-voice-pipeline/internal/adapters/providers"
	turnadapter "port-voice-pipeline/internal/adapters/turn"
	"port-voice-pipeline/internal/config"
	"port-voice-pipeline/internal/domain/voice"
)

func TestNewSessionRunnerUsesDefaultRunnerWhenTurnDisabled(t *testing.T) {
	runner, turnEnabled, err := newSessionRunner(config.Config{}, fakeProviderRuntime())
	if err != nil {
		t.Fatalf("newSessionRunner() error = %v", err)
	}

	if runner == nil {
		t.Fatal("runner = nil")
	}
	if turnEnabled {
		t.Fatal("turnEnabled = true, want false")
	}
}

func TestNewSessionRunnerUsesTurnAwareRunnerWhenTurnEnabled(t *testing.T) {
	runner, turnEnabled, err := newSessionRunner(config.Config{TurnEnabled: true, VADProvider: "noop"}, fakeProviderRuntime())
	if err != nil {
		t.Fatalf("newSessionRunner() error = %v", err)
	}

	if runner == nil {
		t.Fatal("runner = nil")
	}
	if !turnEnabled {
		t.Fatal("turnEnabled = false, want true")
	}
}

func TestNewSessionRunnerReturnsTurnRuntimeError(t *testing.T) {
	_, _, err := newSessionRunner(config.Config{TurnEnabled: true, VADProvider: "silero"}, fakeProviderRuntime())
	if !errors.Is(err, turnadapter.ErrMissingSileroEngine) {
		t.Fatalf("newSessionRunner() error = %v, want %v", err, turnadapter.ErrMissingSileroEngine)
	}
}

func fakeProviderRuntime() providers.Runtime {
	return providers.Runtime{
		STT: fakeSTT{},
		LLM: fakeLLM{},
		TTS: fakeTTS{},
	}
}

type fakeSTT struct{}

func (fakeSTT) Transcribe(context.Context, <-chan voice.PCMFrame) (<-chan voice.Transcript, error) {
	out := make(chan voice.Transcript)
	close(out)
	return out, nil
}

type fakeLLM struct{}

func (fakeLLM) Generate(context.Context, voice.UserUtterance) (voice.AssistantResponse, error) {
	return voice.AssistantResponse{}, nil
}

type fakeTTS struct{}

func (fakeTTS) Synthesize(context.Context, voice.AssistantResponse) (<-chan voice.PCMFrame, error) {
	out := make(chan voice.PCMFrame)
	close(out)
	return out, nil
}
