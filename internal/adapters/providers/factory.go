package providers

import (
	"errors"
	"fmt"
	"strings"

	"port-voice-pipeline/internal/adapters/providers/cartesia"
	"port-voice-pipeline/internal/adapters/providers/deepgram"
	"port-voice-pipeline/internal/adapters/providers/openrouter"
	"port-voice-pipeline/internal/application/ports"
)

var ErrUnsupportedTTSProvider = errors.New("unsupported tts provider")
var ErrMissingProviderConfig = errors.New("missing provider config")

type Runtime struct {
	STT ports.SpeechToText
	LLM ports.LanguageModel
	TTS ports.TextToSpeech
}

type Config struct {
	Deepgram    deepgram.Config
	OpenRouter  openrouter.Config
	Cartesia    cartesia.Config
	TTSProvider string
}

func NewRuntime(cfg Config) (Runtime, error) {
	if err := validateConfig(cfg); err != nil {
		return Runtime{}, err
	}

	tts, err := newTTS(cfg)
	if err != nil {
		return Runtime{}, err
	}

	return Runtime{
		STT: deepgram.New(cfg.Deepgram),
		LLM: openrouter.New(cfg.OpenRouter),
		TTS: tts,
	}, nil
}

func validateConfig(cfg Config) error {
	if cfg.Deepgram.APIKey == "" {
		return fmt.Errorf("%w: DEEPGRAM_API_KEY", ErrMissingProviderConfig)
	}
	if cfg.OpenRouter.APIKey == "" {
		return fmt.Errorf("%w: OPENROUTER_API_KEY", ErrMissingProviderConfig)
	}
	if strings.EqualFold(strings.TrimSpace(cfg.TTSProvider), "cartesia") || strings.TrimSpace(cfg.TTSProvider) == "" {
		if cfg.Cartesia.APIKey == "" {
			return fmt.Errorf("%w: CARTESIA_API_KEY", ErrMissingProviderConfig)
		}
		if cfg.Cartesia.VoiceID == "" {
			return fmt.Errorf("%w: CARTESIA_VOICE_ID", ErrMissingProviderConfig)
		}
	}
	return nil
}

func newTTS(cfg Config) (ports.TextToSpeech, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.TTSProvider)) {
	case "", "cartesia":
		return cartesia.New(cfg.Cartesia), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedTTSProvider, cfg.TTSProvider)
	}
}
