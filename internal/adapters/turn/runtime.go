package turnadapter

import (
	"context"
	"log"
	"time"

	"port-agent-worker/internal/adapters/vad/noop"
	"port-agent-worker/internal/application/session"
	applicationturn "port-agent-worker/internal/application/turn"
	"port-agent-worker/internal/config"
)

const defaultTickInterval = 50 * time.Millisecond

type Logger interface {
	Printf(format string, args ...any)
}

type DecisionLogger struct {
	Logger Logger
}

func (h DecisionLogger) HandleTurnDecision(_ context.Context, decision applicationturn.Decision) error {
	if !decision.BargeIn && !decision.Endpoint {
		return nil
	}

	logger := h.Logger
	if logger == nil {
		logger = log.Default()
	}
	logger.Printf("turn decision barge_in=%t endpoint=%t", decision.BargeIn, decision.Endpoint)
	return nil
}

func NewRuntime(cfg config.Config) session.TurnRuntime {
	controller := applicationturn.NewController(applicationturn.Config{
		StopDelay:        cfg.TurnStopDelay,
		SmartTurnEnabled: cfg.SmartTurnEnabled,
	}, nil)

	return session.TurnRuntime{
		VAD:          noop.Detector{},
		Processor:    applicationturn.NewActivityProcessor(controller),
		Handler:      DecisionLogger{},
		TickInterval: defaultTickInterval,
	}
}
