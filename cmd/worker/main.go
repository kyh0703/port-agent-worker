package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"port-agent-worker/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log.Printf("port-agent-worker starting env=%s tts_provider=%s smart_turn_enabled=%t", cfg.Environment, cfg.TTSProvider, cfg.SmartTurnEnabled)
	log.Print("pion track wiring is pending; worker skeleton is ready")

	<-ctx.Done()
	log.Print("port-agent-worker stopped")
}
