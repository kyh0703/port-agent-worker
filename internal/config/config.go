package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment               string
	DeepgramAPIKey            string
	DeepgramModel             string
	DeepgramLanguage          string
	OpenRouterKey             string
	OpenRouterModel           string
	SystemPrompt              string
	TTSProvider               string
	CartesiaAPIKey            string
	CartesiaVoiceID           string
	CartesiaModelID           string
	ElevenLabsKey             string
	VADProvider               string
	SileroVADThreshold        float64
	SileroVADMinSpeechFrames  int
	SileroVADMinSilenceFrames int
	SmartTurnEnabled          bool
	TurnEnabled               bool
	TurnStopDelay             time.Duration
	RunSession                bool
}

func Load() Config {
	return Config{
		Environment:               env("APP_ENV", "development"),
		DeepgramAPIKey:            os.Getenv("DEEPGRAM_API_KEY"),
		DeepgramModel:             env("DEEPGRAM_MODEL", "nova-3"),
		DeepgramLanguage:          env("DEEPGRAM_LANGUAGE", "ko"),
		OpenRouterKey:             os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterModel:           env("OPENROUTER_MODEL", "google/gemini-2.5-flash-lite"),
		SystemPrompt:              os.Getenv("SYSTEM_PROMPT"),
		TTSProvider:               env("TTS_PROVIDER", "cartesia"),
		CartesiaAPIKey:            os.Getenv("CARTESIA_API_KEY"),
		CartesiaVoiceID:           os.Getenv("CARTESIA_VOICE_ID"),
		CartesiaModelID:           env("CARTESIA_MODEL_ID", "sonic-3.5"),
		ElevenLabsKey:             os.Getenv("ELEVENLABS_API_KEY"),
		VADProvider:               env("VAD_PROVIDER", "noop"),
		SileroVADThreshold:        floatEnv("SILERO_VAD_THRESHOLD", 0.5),
		SileroVADMinSpeechFrames:  intEnv("SILERO_VAD_MIN_SPEECH_FRAMES", 1),
		SileroVADMinSilenceFrames: intEnv("SILERO_VAD_MIN_SILENCE_FRAMES", 3),
		SmartTurnEnabled:          os.Getenv("SMART_TURN_ENABLED") == "true",
		TurnEnabled:               os.Getenv("TURN_ENABLED") == "true",
		TurnStopDelay:             durationEnv("TURN_STOP_DELAY", 700*time.Millisecond),
		RunSession:                os.Getenv("RUN_SESSION") == "true",
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}

func floatEnv(key string, fallback float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func intEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
