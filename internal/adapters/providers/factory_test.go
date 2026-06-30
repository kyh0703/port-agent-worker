package providers

import (
	"errors"
	"testing"

	"port-agent-worker/internal/adapters/providers/cartesia"
	"port-agent-worker/internal/adapters/providers/deepgram"
	"port-agent-worker/internal/adapters/providers/openrouter"
)

func TestNewRuntimeBuildsProviders(t *testing.T) {
	runtime, err := NewRuntime(Config{
		Deepgram: deepgram.Config{
			APIKey:   "deepgram-key",
			Model:    "nova-3",
			Language: "ko",
		},
		OpenRouter: openrouter.Config{
			APIKey: "openrouter-key",
			Model:  "google/gemini-2.5-flash-lite",
		},
		Cartesia: cartesia.Config{
			APIKey:  "cartesia-key",
			VoiceID: "voice-id",
			ModelID: "sonic-3.5",
		},
		TTSProvider: "cartesia",
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
	_, err := NewRuntime(Config{
		Deepgram:    deepgram.Config{APIKey: "deepgram-key"},
		OpenRouter:  openrouter.Config{APIKey: "openrouter-key"},
		TTSProvider: "unknown",
	})
	if !errors.Is(err, ErrUnsupportedTTSProvider) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrUnsupportedTTSProvider)
	}
}

func TestNewRuntimeRequiresProviderConfig(t *testing.T) {
	_, err := NewRuntime(Config{})
	if !errors.Is(err, ErrMissingProviderConfig) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrMissingProviderConfig)
	}

	_, err = NewRuntime(Config{
		Deepgram:    deepgram.Config{APIKey: "deepgram-key"},
		OpenRouter:  openrouter.Config{APIKey: "openrouter-key"},
		TTSProvider: "cartesia",
		Cartesia:    cartesia.Config{APIKey: "cartesia-key"},
	})
	if !errors.Is(err, ErrMissingProviderConfig) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrMissingProviderConfig)
	}
}
