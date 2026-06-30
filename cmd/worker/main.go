package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"port-agent-worker/internal/adapters/providers"
	"port-agent-worker/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log.Printf("port-agent-worker starting env=%s tts_provider=%s smart_turn_enabled=%t", cfg.Environment, cfg.TTSProvider, cfg.SmartTurnEnabled)

	runtime, err := providers.NewRuntime(cfg)
	if err != nil {
		log.Printf("provider wiring failed: %v", err)
		os.Exit(1)
	}
	if runtime.STT == nil || runtime.LLM == nil || runtime.TTS == nil {
		log.Print("provider wiring failed: nil provider")
		os.Exit(1)
	}

	log.Print("provider wiring ready")
	log.Print("pion track wiring is pending")

	<-ctx.Done()
	log.Print("port-agent-worker stopped")
}
