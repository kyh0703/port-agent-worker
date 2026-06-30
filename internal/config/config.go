package config

import "os"

type Config struct {
	Environment      string
	DeepgramAPIKey   string
	DeepgramModel    string
	DeepgramLanguage string
	OpenRouterKey    string
	OpenRouterModel  string
	SystemPrompt     string
	TTSProvider      string
	CartesiaAPIKey   string
	CartesiaVoiceID  string
	CartesiaModelID  string
	ElevenLabsKey    string
	SmartTurnEnabled bool
	RunSession       bool
}

func Load() Config {
	return Config{
		Environment:      env("APP_ENV", "development"),
		DeepgramAPIKey:   os.Getenv("DEEPGRAM_API_KEY"),
		DeepgramModel:    env("DEEPGRAM_MODEL", "nova-3"),
		DeepgramLanguage: env("DEEPGRAM_LANGUAGE", "ko"),
		OpenRouterKey:    os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterModel:  env("OPENROUTER_MODEL", "google/gemini-2.5-flash-lite"),
		SystemPrompt:     os.Getenv("SYSTEM_PROMPT"),
		TTSProvider:      env("TTS_PROVIDER", "cartesia"),
		CartesiaAPIKey:   os.Getenv("CARTESIA_API_KEY"),
		CartesiaVoiceID:  os.Getenv("CARTESIA_VOICE_ID"),
		CartesiaModelID:  env("CARTESIA_MODEL_ID", "sonic-3.5"),
		ElevenLabsKey:    os.Getenv("ELEVENLABS_API_KEY"),
		SmartTurnEnabled: os.Getenv("SMART_TURN_ENABLED") == "true",
		RunSession:       os.Getenv("RUN_SESSION") == "true",
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
