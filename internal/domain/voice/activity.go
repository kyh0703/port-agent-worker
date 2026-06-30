package voice

import (
	"errors"
	"time"
)

type SpeechActivityKind int

const (
	SpeechStarted SpeechActivityKind = iota + 1
	SpeechStopped
)

var ErrInvalidSpeechActivityEvent = errors.New("invalid speech activity event")

type SpeechActivityEvent struct {
	Kind SpeechActivityKind
	At   time.Time
}

func NewSpeechActivityEvent(kind SpeechActivityKind, at time.Time) (SpeechActivityEvent, error) {
	if !kind.Valid() || at.IsZero() {
		return SpeechActivityEvent{}, ErrInvalidSpeechActivityEvent
	}

	return SpeechActivityEvent{
		Kind: kind,
		At:   at,
	}, nil
}

func (k SpeechActivityKind) Valid() bool {
	return k == SpeechStarted || k == SpeechStopped
}
