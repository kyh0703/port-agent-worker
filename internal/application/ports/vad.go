package ports

import (
	"context"

	"port-voice-pipeline/internal/domain/voice"
)

type VoiceActivityDetector interface {
	DetectSpeech(ctx context.Context, frames <-chan voice.PCMFrame) (<-chan voice.SpeechActivityEvent, error)
}
