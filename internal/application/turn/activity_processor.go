package turn

import (
	"context"
	"errors"
	"time"

	"port-voice-pipeline/internal/domain/voice"
)

var ErrControllerRequired = errors.New("turn controller required")

type ActivityProcessor struct {
	controller *Controller
}

func NewActivityProcessor(controller *Controller) *ActivityProcessor {
	return &ActivityProcessor{controller: controller}
}

func (p *ActivityProcessor) Handle(ctx context.Context, event voice.SpeechActivityEvent) (Decision, error) {
	if p.controller == nil {
		return Decision{}, ErrControllerRequired
	}
	if !event.Kind.Valid() || event.At.IsZero() {
		return Decision{}, voice.ErrInvalidSpeechActivityEvent
	}

	switch event.Kind {
	case voice.SpeechStarted:
		return p.controller.UserSpeechStarted(event.At), nil
	case voice.SpeechStopped:
		p.controller.UserSpeechStopped(event.At)
		return p.Tick(ctx, event.At)
	default:
		return Decision{}, voice.ErrInvalidSpeechActivityEvent
	}
}

func (p *ActivityProcessor) Tick(ctx context.Context, now time.Time) (Decision, error) {
	if p.controller == nil {
		return Decision{}, ErrControllerRequired
	}

	endpoint, err := p.controller.ShouldEndpoint(ctx, now)
	if err != nil {
		return Decision{}, err
	}

	return Decision{Endpoint: endpoint}, nil
}
