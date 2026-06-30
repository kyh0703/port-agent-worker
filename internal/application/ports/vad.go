package ports

import (
	"context"

	"port-agent-worker/internal/domain/voice"
)

type VoiceActivityDetector interface {
	DetectSpeech(ctx context.Context, frames <-chan voice.PCMFrame) (<-chan voice.SpeechActivityEvent, error)
}
