package main

import (
	"context"
	"testing"

	"port-agent-worker/internal/adapters/providers"
	"port-agent-worker/internal/config"
	"port-agent-worker/internal/domain/voice"
)

func TestNewSessionRunnerUsesDefaultRunnerWhenTurnDisabled(t *testing.T) {
	runner, turnEnabled := newSessionRunner(config.Config{}, fakeProviderRuntime())

	if runner == nil {
		t.Fatal("runner = nil")
	}
	if turnEnabled {
		t.Fatal("turnEnabled = true, want false")
	}
}

func TestNewSessionRunnerUsesTurnAwareRunnerWhenTurnEnabled(t *testing.T) {
	runner, turnEnabled := newSessionRunner(config.Config{TurnEnabled: true}, fakeProviderRuntime())

	if runner == nil {
		t.Fatal("runner = nil")
	}
	if !turnEnabled {
		t.Fatal("turnEnabled = false, want true")
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
