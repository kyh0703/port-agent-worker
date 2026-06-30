package turnadapter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"port-agent-worker/internal/adapters/vad/noop"
	"port-agent-worker/internal/adapters/vad/silero"
	"port-agent-worker/internal/application/ports"
	"port-agent-worker/internal/application/session"
	applicationturn "port-agent-worker/internal/application/turn"
	"port-agent-worker/internal/config"
)

const defaultTickInterval = 50 * time.Millisecond

var ErrUnsupportedVADProvider = errors.New("unsupported vad provider")
var ErrMissingSileroEngine = errors.New("missing silero engine")

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

type RuntimeOptions struct {
	SileroEngine silero.Engine
}

func NewRuntime(cfg config.Config) (session.TurnRuntime, error) {
	return NewRuntimeWithOptions(cfg, RuntimeOptions{})
}

func NewRuntimeWithOptions(cfg config.Config, options RuntimeOptions) (session.TurnRuntime, error) {
	controller := applicationturn.NewController(applicationturn.Config{
		StopDelay:        cfg.TurnStopDelay,
		SmartTurnEnabled: cfg.SmartTurnEnabled,
	}, nil)

	vad, err := newVAD(cfg, options)
	if err != nil {
		return session.TurnRuntime{}, err
	}

	return session.TurnRuntime{
		VAD:          vad,
		Processor:    applicationturn.NewActivityProcessor(controller),
		Handler:      DecisionLogger{},
		TickInterval: defaultTickInterval,
	}, nil
}

func newVAD(cfg config.Config, options RuntimeOptions) (ports.VoiceActivityDetector, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.VADProvider)) {
	case "", "noop":
		return noop.Detector{}, nil
	case "silero":
		if options.SileroEngine == nil {
			return nil, ErrMissingSileroEngine
		}
		detector, err := silero.New(options.SileroEngine, silero.Config{
			Threshold:        cfg.SileroVADThreshold,
			MinSpeechFrames:  cfg.SileroVADMinSpeechFrames,
			MinSilenceFrames: cfg.SileroVADMinSilenceFrames,
		})
		if err != nil {
			return nil, err
		}
		return detector, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedVADProvider, cfg.VADProvider)
	}
}
