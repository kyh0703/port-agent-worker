package session

import (
	"context"
	"errors"
	"time"
)

const defaultIdleDelay = 50 * time.Millisecond

type TurnExecutor interface {
	RunTurn(ctx context.Context) error
}

type Runner struct {
	executor  TurnExecutor
	idleDelay time.Duration
}

func NewRunner(executor TurnExecutor) *Runner {
	return &Runner{
		executor:  executor,
		idleDelay: defaultIdleDelay,
	}
}

func NewRunnerWithIdleDelay(executor TurnExecutor, idleDelay time.Duration) *Runner {
	if idleDelay <= 0 {
		idleDelay = defaultIdleDelay
	}
	return &Runner{
		executor:  executor,
		idleDelay: idleDelay,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := r.executor.RunTurn(ctx)
		if err == nil {
			continue
		}
		if errors.Is(err, ErrNoFinalTranscript) {
			if waitErr := sleep(ctx, r.idleDelay); waitErr != nil {
				return waitErr
			}
			continue
		}
		return err
	}
}

func sleep(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
