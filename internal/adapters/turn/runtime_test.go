package turnadapter

import (
	"context"
	"errors"
	"testing"
	"time"

	"port-voice-pipeline/internal/adapters/vad/noop"
	"port-voice-pipeline/internal/adapters/vad/silero"
	"port-voice-pipeline/internal/application/turn"
	"port-voice-pipeline/internal/config"
	"port-voice-pipeline/internal/domain/voice"
)

func TestNewRuntimeBuildsTurnRuntime(t *testing.T) {
	runtime, err := NewRuntime(config.Config{
		SmartTurnEnabled: true,
		TurnStopDelay:    time.Second,
		VADProvider:      "noop",
	})
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}

	if runtime.VAD == nil {
		t.Fatal("VAD = nil")
	}
	if _, ok := runtime.VAD.(noop.Detector); !ok {
		t.Fatalf("VAD type = %T, want noop.Detector", runtime.VAD)
	}
	if runtime.Processor == nil {
		t.Fatal("Processor = nil")
	}
	if runtime.Handler == nil {
		t.Fatal("Handler = nil")
	}
	if runtime.TickInterval <= 0 {
		t.Fatalf("TickInterval = %s, want positive duration", runtime.TickInterval)
	}
}

func TestNewRuntimeRejectsUnknownVADProvider(t *testing.T) {
	_, err := NewRuntime(config.Config{VADProvider: "unknown"})
	if !errors.Is(err, ErrUnsupportedVADProvider) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrUnsupportedVADProvider)
	}
}

func TestNewRuntimeRequiresSileroEngine(t *testing.T) {
	_, err := NewRuntime(config.Config{VADProvider: "silero"})
	if !errors.Is(err, ErrMissingSileroEngine) {
		t.Fatalf("NewRuntime() error = %v, want %v", err, ErrMissingSileroEngine)
	}
}

func TestNewRuntimeWithOptionsBuildsSileroRuntime(t *testing.T) {
	runtime, err := NewRuntimeWithOptions(
		config.Config{
			VADProvider:               "silero",
			SileroVADThreshold:        0.7,
			SileroVADMinSpeechFrames:  2,
			SileroVADMinSilenceFrames: 4,
		},
		RuntimeOptions{SileroEngine: fixedSileroEngine{}},
	)
	if err != nil {
		t.Fatalf("NewRuntimeWithOptions() error = %v", err)
	}
	if _, ok := runtime.VAD.(*silero.Detector); !ok {
		t.Fatalf("VAD type = %T, want *silero.Detector", runtime.VAD)
	}
}

func TestDecisionLoggerLogsOnlyNonEmptyDecision(t *testing.T) {
	logger := &recordingLogger{}
	handler := DecisionLogger{Logger: logger}

	if err := handler.HandleTurnDecision(context.Background(), turn.Decision{}); err != nil {
		t.Fatalf("HandleTurnDecision() error = %v", err)
	}
	if logger.calls != 0 {
		t.Fatalf("logger calls = %d, want 0", logger.calls)
	}

	if err := handler.HandleTurnDecision(context.Background(), turn.Decision{BargeIn: true}); err != nil {
		t.Fatalf("HandleTurnDecision() error = %v", err)
	}
	if logger.calls != 1 {
		t.Fatalf("logger calls = %d, want 1", logger.calls)
	}
}

type recordingLogger struct {
	calls int
}

func (l *recordingLogger) Printf(string, ...any) {
	l.calls++
}

type fixedSileroEngine struct{}

func (fixedSileroEngine) SpeechProbability(context.Context, voice.PCMFrame) (float64, error) {
	return 0, nil
}
