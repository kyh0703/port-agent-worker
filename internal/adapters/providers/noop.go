package providers

import (
	"context"
	"errors"

	"port-voice-pipeline/internal/domain/voice"
)

var ErrProviderNotConfigured = errors.New("provider is not configured")

type NoopSTT struct{}

func (NoopSTT) Transcribe(context.Context, <-chan voice.PCMFrame) (<-chan voice.Transcript, error) {
	return nil, ErrProviderNotConfigured
}

type NoopLLM struct{}

func (NoopLLM) Generate(context.Context, voice.UserUtterance) (voice.AssistantResponse, error) {
	return voice.AssistantResponse{}, ErrProviderNotConfigured
}

type NoopTTS struct{}

func (NoopTTS) Synthesize(context.Context, voice.AssistantResponse) (<-chan voice.PCMFrame, error) {
	return nil, ErrProviderNotConfigured
}
