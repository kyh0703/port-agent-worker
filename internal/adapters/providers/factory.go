package providers

import (
	"errors"
	"fmt"
	"strings"

	"port-agent-worker/internal/adapters/providers/cartesia"
	"port-agent-worker/internal/adapters/providers/deepgram"
	"port-agent-worker/internal/adapters/providers/openrouter"
	"port-agent-worker/internal/application/ports"
	"port-agent-worker/internal/config"
)

var ErrUnsupportedTTSProvider = errors.New("unsupported tts provider")
var ErrMissingProviderConfig = errors.New("missing provider config")

type Runtime struct {
	STT ports.SpeechToText
	LLM ports.LanguageModel
	TTS ports.TextToSpeech
}

func NewRuntime(cfg config.Config) (Runtime, error) {
	if err := validateConfig(cfg); err != nil {
		return Runtime{}, err
	}

	tts, err := newTTS(cfg)
	if err != nil {
		return Runtime{}, err
	}

	return Runtime{
		STT: deepgram.New(deepgram.Config{
			APIKey:         cfg.DeepgramAPIKey,
			Model:          cfg.DeepgramModel,
			Language:       cfg.DeepgramLanguage,
			InterimResults: true,
			SmartFormat:    true,
		}),
		LLM: openrouter.New(openrouter.Config{
			APIKey:       cfg.OpenRouterKey,
			Model:        cfg.OpenRouterModel,
			SystemPrompt: cfg.SystemPrompt,
			AppTitle:     "port-agent-worker",
		}),
		TTS: tts,
	}, nil
}

func validateConfig(cfg config.Config) error {
	if cfg.DeepgramAPIKey == "" {
		return fmt.Errorf("%w: DEEPGRAM_API_KEY", ErrMissingProviderConfig)
	}
	if cfg.OpenRouterKey == "" {
		return fmt.Errorf("%w: OPENROUTER_API_KEY", ErrMissingProviderConfig)
	}
	if strings.EqualFold(strings.TrimSpace(cfg.TTSProvider), "cartesia") || strings.TrimSpace(cfg.TTSProvider) == "" {
		if cfg.CartesiaAPIKey == "" {
			return fmt.Errorf("%w: CARTESIA_API_KEY", ErrMissingProviderConfig)
		}
		if cfg.CartesiaVoiceID == "" {
			return fmt.Errorf("%w: CARTESIA_VOICE_ID", ErrMissingProviderConfig)
		}
	}
	return nil
}

func newTTS(cfg config.Config) (ports.TextToSpeech, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.TTSProvider)) {
	case "", "cartesia":
		return cartesia.New(cartesia.Config{
			APIKey:  cfg.CartesiaAPIKey,
			ModelID: cfg.CartesiaModelID,
			VoiceID: cfg.CartesiaVoiceID,
		}), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedTTSProvider, cfg.TTSProvider)
	}
}
