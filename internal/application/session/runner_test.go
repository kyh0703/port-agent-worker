package session

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunnerRepeatsTurnsUntilContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	executor := &fakeTurnExecutor{
		onCall: func(call int) error {
			if call == 3 {
				cancel()
			}
			return nil
		},
	}

	err := NewRunnerWithIdleDelay(executor, time.Nanosecond).Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want %v", err, context.Canceled)
	}
	if executor.calls != 3 {
		t.Fatalf("calls = %d, want 3", executor.calls)
	}
}

func TestRunnerTreatsNoFinalTranscriptAsIdle(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	executor := &fakeTurnExecutor{
		onCall: func(call int) error {
			if call == 2 {
				cancel()
			}
			return ErrNoFinalTranscript
		},
	}

	err := NewRunnerWithIdleDelay(executor, time.Nanosecond).Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want %v", err, context.Canceled)
	}
	if executor.calls != 2 {
		t.Fatalf("calls = %d, want 2", executor.calls)
	}
}

func TestRunnerReturnsFatalError(t *testing.T) {
	fatalErr := errors.New("fatal")
	executor := &fakeTurnExecutor{
		onCall: func(int) error {
			return fatalErr
		},
	}

	err := NewRunnerWithIdleDelay(executor, time.Nanosecond).Run(context.Background())
	if !errors.Is(err, fatalErr) {
		t.Fatalf("Run() error = %v, want %v", err, fatalErr)
	}
	if executor.calls != 1 {
		t.Fatalf("calls = %d, want 1", executor.calls)
	}
}

type fakeTurnExecutor struct {
	calls  int
	onCall func(call int) error
}

func (f *fakeTurnExecutor) RunTurn(context.Context) error {
	f.calls++
	return f.onCall(f.calls)
}
