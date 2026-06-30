package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("DEEPGRAM_MODEL", "")
	t.Setenv("DEEPGRAM_LANGUAGE", "")
	t.Setenv("OPENROUTER_MODEL", "")
	t.Setenv("TTS_PROVIDER", "")
	t.Setenv("CARTESIA_MODEL_ID", "")
	t.Setenv("VAD_PROVIDER", "")
	t.Setenv("SILERO_VAD_THRESHOLD", "")
	t.Setenv("SILERO_VAD_MIN_SPEECH_FRAMES", "")
	t.Setenv("SILERO_VAD_MIN_SILENCE_FRAMES", "")
	t.Setenv("TURN_ENABLED", "")
	t.Setenv("TURN_STOP_DELAY", "")

	cfg := Load()

	if cfg.Environment != "development" {
		t.Fatalf("Environment = %q, want development", cfg.Environment)
	}
	if cfg.DeepgramModel != "nova-3" {
		t.Fatalf("DeepgramModel = %q, want nova-3", cfg.DeepgramModel)
	}
	if cfg.DeepgramLanguage != "ko" {
		t.Fatalf("DeepgramLanguage = %q, want ko", cfg.DeepgramLanguage)
	}
	if cfg.OpenRouterModel != "google/gemini-2.5-flash-lite" {
		t.Fatalf("OpenRouterModel = %q", cfg.OpenRouterModel)
	}
	if cfg.TTSProvider != "cartesia" {
		t.Fatalf("TTSProvider = %q, want cartesia", cfg.TTSProvider)
	}
	if cfg.CartesiaModelID != "sonic-3.5" {
		t.Fatalf("CartesiaModelID = %q, want sonic-3.5", cfg.CartesiaModelID)
	}
	if cfg.VADProvider != "noop" {
		t.Fatalf("VADProvider = %q, want noop", cfg.VADProvider)
	}
	if cfg.SileroVADThreshold != 0.5 {
		t.Fatalf("SileroVADThreshold = %f, want 0.5", cfg.SileroVADThreshold)
	}
	if cfg.SileroVADMinSpeechFrames != 1 {
		t.Fatalf("SileroVADMinSpeechFrames = %d, want 1", cfg.SileroVADMinSpeechFrames)
	}
	if cfg.SileroVADMinSilenceFrames != 3 {
		t.Fatalf("SileroVADMinSilenceFrames = %d, want 3", cfg.SileroVADMinSilenceFrames)
	}
	if cfg.TurnEnabled {
		t.Fatal("TurnEnabled = true, want false")
	}
	if cfg.TurnStopDelay.String() != "700ms" {
		t.Fatalf("TurnStopDelay = %s, want 700ms", cfg.TurnStopDelay)
	}
}

func TestLoadProviderSecrets(t *testing.T) {
	t.Setenv("DEEPGRAM_API_KEY", "deepgram-key")
	t.Setenv("OPENROUTER_API_KEY", "openrouter-key")
	t.Setenv("CARTESIA_API_KEY", "cartesia-key")
	t.Setenv("CARTESIA_VOICE_ID", "voice-id")
	t.Setenv("SYSTEM_PROMPT", "be brief")
	t.Setenv("SMART_TURN_ENABLED", "true")
	t.Setenv("VAD_PROVIDER", "silero")
	t.Setenv("SILERO_VAD_THRESHOLD", "0.7")
	t.Setenv("SILERO_VAD_MIN_SPEECH_FRAMES", "2")
	t.Setenv("SILERO_VAD_MIN_SILENCE_FRAMES", "4")
	t.Setenv("TURN_ENABLED", "true")
	t.Setenv("TURN_STOP_DELAY", "2s")
	t.Setenv("RUN_SESSION", "true")

	cfg := Load()

	if cfg.DeepgramAPIKey != "deepgram-key" {
		t.Fatalf("DeepgramAPIKey not loaded")
	}
	if cfg.OpenRouterKey != "openrouter-key" {
		t.Fatalf("OpenRouterKey not loaded")
	}
	if cfg.CartesiaAPIKey != "cartesia-key" {
		t.Fatalf("CartesiaAPIKey not loaded")
	}
	if cfg.CartesiaVoiceID != "voice-id" {
		t.Fatalf("CartesiaVoiceID = %q, want voice-id", cfg.CartesiaVoiceID)
	}
	if cfg.SystemPrompt != "be brief" {
		t.Fatalf("SystemPrompt = %q, want be brief", cfg.SystemPrompt)
	}
	if !cfg.SmartTurnEnabled {
		t.Fatal("SmartTurnEnabled = false, want true")
	}
	if cfg.VADProvider != "silero" {
		t.Fatalf("VADProvider = %q, want silero", cfg.VADProvider)
	}
	if cfg.SileroVADThreshold != 0.7 {
		t.Fatalf("SileroVADThreshold = %f, want 0.7", cfg.SileroVADThreshold)
	}
	if cfg.SileroVADMinSpeechFrames != 2 {
		t.Fatalf("SileroVADMinSpeechFrames = %d, want 2", cfg.SileroVADMinSpeechFrames)
	}
	if cfg.SileroVADMinSilenceFrames != 4 {
		t.Fatalf("SileroVADMinSilenceFrames = %d, want 4", cfg.SileroVADMinSilenceFrames)
	}
	if !cfg.TurnEnabled {
		t.Fatal("TurnEnabled = false, want true")
	}
	if cfg.TurnStopDelay.String() != "2s" {
		t.Fatalf("TurnStopDelay = %s, want 2s", cfg.TurnStopDelay)
	}
	if !cfg.RunSession {
		t.Fatal("RunSession = false, want true")
	}
}
