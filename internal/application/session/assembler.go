package session

import (
	"port-voice-pipeline/internal/application/ports"
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

func NewTurnAwareRunnerFromRuntime(providers ProviderRuntime, audio AudioRuntime, turnRuntime TurnRuntime) *Runner {
	orchestrator := NewTurnAwareOrchestrator(
		audio.Ingress,
		audio.Egress,
		providers.STT,
		providers.LLM,
		providers.TTS,
		turnRuntime,
	)
	return NewRunner(orchestrator)
}
