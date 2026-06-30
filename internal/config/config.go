package config

import "os"

type Config struct {
	Environment      string
	DeepgramAPIKey   string
	OpenRouterKey    string
	TTSProvider      string
	CartesiaAPIKey   string
	ElevenLabsKey    string
	SmartTurnEnabled bool
}

func Load() Config {
	return Config{
		Environment:      env("APP_ENV", "development"),
		DeepgramAPIKey:   os.Getenv("DEEPGRAM_API_KEY"),
		OpenRouterKey:    os.Getenv("OPENROUTER_API_KEY"),
		TTSProvider:      env("TTS_PROVIDER", "cartesia"),
		CartesiaAPIKey:   os.Getenv("CARTESIA_API_KEY"),
		ElevenLabsKey:    os.Getenv("ELEVENLABS_API_KEY"),
		SmartTurnEnabled: os.Getenv("SMART_TURN_ENABLED") == "true",
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
