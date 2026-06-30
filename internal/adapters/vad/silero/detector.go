package silero

import (
	"context"
	"errors"
	"time"

	"port-agent-worker/internal/domain/voice"
)

const (
	defaultThreshold        = 0.5
	defaultMinSpeechFrames  = 1
	defaultMinSilenceFrames = 3
)

var ErrEngineRequired = errors.New("silero engine required")
var ErrInvalidConfig = errors.New("invalid silero vad config")

type Engine interface {
	SpeechProbability(ctx context.Context, frame voice.PCMFrame) (float64, error)
}

type Config struct {
	Threshold        float64
	MinSpeechFrames  int
	MinSilenceFrames int
}

type Detector struct {
	engine           Engine
	threshold        float64
	minSpeechFrames  int
	minSilenceFrames int
	now              func() time.Time
}

func New(engine Engine, config Config) (*Detector, error) {
	if engine == nil {
		return nil, ErrEngineRequired
	}

	config = normalizeConfig(config)
	if config.Threshold <= 0 || config.Threshold >= 1 || config.MinSpeechFrames <= 0 || config.MinSilenceFrames <= 0 {
		return nil, ErrInvalidConfig
	}

	return &Detector{
		engine:           engine,
		threshold:        config.Threshold,
		minSpeechFrames:  config.MinSpeechFrames,
		minSilenceFrames: config.MinSilenceFrames,
		now:              time.Now,
	}, nil
}

func (d *Detector) DetectSpeech(ctx context.Context, frames <-chan voice.PCMFrame) (<-chan voice.SpeechActivityEvent, error) {
	events := make(chan voice.SpeechActivityEvent)
	go func() {
		defer close(events)
		state := detectorState{}

		for {
			select {
			case <-ctx.Done():
				return
			case frame, ok := <-frames:
				if !ok {
					return
				}
				probability, err := d.engine.SpeechProbability(ctx, frame)
				if err != nil {
					return
				}
				event, emit := d.applyProbability(&state, probability)
				if !emit {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case events <- event:
				}
			}
		}
	}()
	return events, nil
}

type detectorState struct {
	speechActive  bool
	speechFrames  int
	silenceFrames int
}

func (d *Detector) applyProbability(state *detectorState, probability float64) (voice.SpeechActivityEvent, bool) {
	if probability >= d.threshold {
		state.speechFrames++
		state.silenceFrames = 0

		if !state.speechActive && state.speechFrames >= d.minSpeechFrames {
			state.speechActive = true
			return voice.SpeechActivityEvent{Kind: voice.SpeechStarted, At: d.now()}, true
		}
		return voice.SpeechActivityEvent{}, false
	}

	state.silenceFrames++
	state.speechFrames = 0

	if state.speechActive && state.silenceFrames >= d.minSilenceFrames {
		state.speechActive = false
		return voice.SpeechActivityEvent{Kind: voice.SpeechStopped, At: d.now()}, true
	}
	return voice.SpeechActivityEvent{}, false
}

func normalizeConfig(config Config) Config {
	if config.Threshold == 0 {
		config.Threshold = defaultThreshold
	}
	if config.MinSpeechFrames == 0 {
		config.MinSpeechFrames = defaultMinSpeechFrames
	}
	if config.MinSilenceFrames == 0 {
		config.MinSilenceFrames = defaultMinSilenceFrames
	}
	return config
}
