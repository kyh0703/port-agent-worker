package pion

import (
	"port-agent-worker/internal/application/session"
)

func NewRunner(config Config, providers session.ProviderRuntime) (*session.Runner, error) {
	audio, err := NewRuntime(config)
	if err != nil {
		return nil, err
	}

	return session.NewRunnerFromRuntime(providers, audio), nil
}
