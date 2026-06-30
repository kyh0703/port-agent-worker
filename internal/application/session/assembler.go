package session

import (
	"port-agent-worker/internal/application/ports"
)

type ProviderRuntime struct {
	STT ports.SpeechToText
	LLM ports.LanguageModel
	TTS ports.TextToSpeech
}

type AudioRuntime struct {
	Ingress ports.AudioIngress
	Egress  ports.AudioEgress
}

func NewRunnerFromRuntime(providers ProviderRuntime, audio AudioRuntime) *Runner {
	orchestrator := NewOrchestrator(
		audio.Ingress,
		audio.Egress,
		providers.STT,
		providers.LLM,
		providers.TTS,
	)
	return NewRunner(orchestrator)
}
