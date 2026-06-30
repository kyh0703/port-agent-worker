package noop

import (
	"context"

	"port-agent-worker/internal/domain/voice"
)

type Detector struct{}

func (Detector) DetectSpeech(ctx context.Context, frames <-chan voice.PCMFrame) (<-chan voice.SpeechActivityEvent, error) {
	events := make(chan voice.SpeechActivityEvent)
	go func() {
		defer close(events)
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-frames:
				if !ok {
					return
				}
			}
		}
	}()
	return events, nil
}
