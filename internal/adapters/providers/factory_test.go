package providers

import (
	"errors"
	"testing"

	"port-agent-worker/internal/config"
)

func TestNewRuntimeBuildsProviders(t *testing.T) {
	runtime, err := NewRuntime(config.Config{
		DeepgramAPIKey:   "deepgram-key",
		DeepgramModel:    "nova-3",
		DeepgramLanguage: "ko",
		OpenRouterKey:    "openrouter-key",
		OpenRouterModel:  "google/gemini-2.5-flash-lite",
		CartesiaAPIKey:   "cartesia-key",
		CartesiaVoiceID:  "voice-id",
		CartesiaModelID:  "sonic-3.5",
		TTSProvider:      "cartesia",
	})
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	if runtime.STT == nil {
		t.Fatal("STT = nil")
	}
	if runtime.LLM == nil {
		t.Fatal("LLM = nil")
	}
	if runtime.TTS == nil {
		t.Fatal("TTS = nil")
	}
}

func TestNewRuntimeRejectsUnsupportedTTSProvider(t *testing.T) {
	_, err := NewRuntime(config.Config{
		DeepgramAPIKey: "deepgram-key",
		OpenRouterKey:  "openrouter-key",
		TTSProvider:    "unknown",
	})
	if !errors.Is(err, ErrUnsupportedTTSProvider) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrUnsupportedTTSProvider)
	}
}

func TestNewRuntimeRequiresProviderConfig(t *testing.T) {
	_, err := NewRuntime(config.Config{})
	if !errors.Is(err, ErrMissingProviderConfig) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrMissingProviderConfig)
	}

	_, err = NewRuntime(config.Config{
		DeepgramAPIKey: "deepgram-key",
		OpenRouterKey:  "openrouter-key",
		TTSProvider:    "cartesia",
		CartesiaAPIKey: "cartesia-key",
	})
	if !errors.Is(err, ErrMissingProviderConfig) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrMissingProviderConfig)
	}
}
