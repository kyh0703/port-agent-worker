package turnadapter

import (
	"context"
	"testing"
	"time"

	"port-agent-worker/internal/application/turn"
	"port-agent-worker/internal/config"
)

func TestNewRuntimeBuildsTurnRuntime(t *testing.T) {
	runtime := NewRuntime(config.Config{
		SmartTurnEnabled: true,
		TurnStopDelay:    time.Second,
	})

	if runtime.VAD == nil {
		t.Fatal("VAD = nil")
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
