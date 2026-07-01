package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"port-voice-pipeline/internal/adapters/media/pending"
	"port-voice-pipeline/internal/adapters/providers"
	"port-voice-pipeline/internal/adapters/providers/cartesia"
	"port-voice-pipeline/internal/adapters/providers/deepgram"
	"port-voice-pipeline/internal/adapters/providers/openrouter"
	turnadapter "port-voice-pipeline/internal/adapters/turn"
	"port-voice-pipeline/internal/application/session"
	"port-voice-pipeline/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log.Printf("port-voice-pipeline starting env=%s tts_provider=%s smart_turn_enabled=%t", cfg.Environment, cfg.TTSProvider, cfg.SmartTurnEnabled)

	runtime, err := newProviderRuntime(cfg)
	if err != nil {
		log.Printf("provider wiring failed: %v", err)
		os.Exit(1)
	}
	if runtime.STT == nil || runtime.LLM == nil || runtime.TTS == nil {
		log.Print("provider wiring failed: nil provider")
		os.Exit(1)
	}

	log.Print("provider wiring ready")
	runner, turnEnabled, err := newSessionRunner(cfg, runtime)
	if err != nil {
		log.Printf("turn runtime wiring failed: %v", err)
		os.Exit(1)
	}
	if turnEnabled {
		log.Print("turn runtime wiring ready")
	}

	if cfg.RunSession {
		log.Print("session runner starting")
		if err := runner.Run(ctx); err != nil {
			log.Printf("session runner stopped: %v", err)
			os.Exit(1)
		}
		return
	}

	log.Print("session runner ready")
	log.Print("media wiring is pending; set RUN_SESSION=true after Pion tracks are configured")

	<-ctx.Done()
	log.Print("port-voice-pipeline stopped")
}

func newProviderRuntime(cfg config.Config) (providers.Runtime, error) {
	return providers.NewRuntime(providers.Config{
		Deepgram: deepgram.Config{
			APIKey:         cfg.DeepgramAPIKey,
			Model:          cfg.DeepgramModel,
			Language:       cfg.DeepgramLanguage,
			InterimResults: true,
			SmartFormat:    true,
		},
		OpenRouter: openrouter.Config{
			APIKey:       cfg.OpenRouterKey,
			Model:        cfg.OpenRouterModel,
			SystemPrompt: cfg.SystemPrompt,
			AppTitle:     "port-voice-pipeline",
		},
		Cartesia: cartesia.Config{
			APIKey:  cfg.CartesiaAPIKey,
			ModelID: cfg.CartesiaModelID,
			VoiceID: cfg.CartesiaVoiceID,
		},
		TTSProvider: cfg.TTSProvider,
	})
}

func newSessionRunner(cfg config.Config, providerRuntime providers.Runtime) (*session.Runner, bool, error) {
	providers := session.ProviderRuntime{
		STT: providerRuntime.STT,
		LLM: providerRuntime.LLM,
		TTS: providerRuntime.TTS,
	}
	audio := session.AudioRuntime{
		Ingress: pending.Ingress{},
		Egress:  pending.Egress{},
	}

	if cfg.TurnEnabled {
		turnRuntime, err := turnadapter.NewRuntime(cfg)
		if err != nil {
			return nil, false, err
		}
		return session.NewTurnAwareRunnerFromRuntime(providers, audio, turnRuntime), true, nil
	}

	return session.NewRunnerFromRuntime(providers, audio), false, nil
}
