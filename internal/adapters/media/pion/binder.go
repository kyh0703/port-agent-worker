package pion

import (
	"port-voice-pipeline/internal/application/session"
)

func NewRunner(config Config, providers session.ProviderRuntime) (*session.Runner, error) {
	audio, err := NewRuntime(config)
	if err != nil {
		return nil, err
	}

	return session.NewRunnerFromRuntime(providers, audio), nil
}

func NewTurnAwareRunner(config Config, providers session.ProviderRuntime, turnRuntime session.TurnRuntime) (*session.Runner, error) {
	audio, err := NewRuntime(config)
	if err != nil {
		return nil, err
	}

	return session.NewTurnAwareRunnerFromRuntime(providers, audio, turnRuntime), nil
}
