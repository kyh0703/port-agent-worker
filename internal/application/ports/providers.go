package ports

import (
	"context"

	"port-agent-worker/internal/domain/voice"
)

type SpeechToText interface {
	Transcribe(ctx context.Context, audio <-chan voice.PCMFrame) (<-chan voice.Transcript, error)
}

type LanguageModel interface {
	Generate(ctx context.Context, utterance voice.UserUtterance) (voice.AssistantResponse, error)
}

type TextToSpeech interface {
	Synthesize(ctx context.Context, response voice.AssistantResponse) (<-chan voice.PCMFrame, error)
}
