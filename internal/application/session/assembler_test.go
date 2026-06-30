package session

import (
	"testing"

	"port-agent-worker/internal/domain/voice"
)

func TestNewRunnerFromRuntimeBuildsRunner(t *testing.T) {
	runner := NewRunnerFromRuntime(
		ProviderRuntime{
			STT: &fakeSTT{transcripts: []voice.Transcript{{Text: "hello", IsFinal: true}}},
			LLM: &fakeLLM{response: voice.AssistantResponse{Text: "hi"}},
			TTS: &fakeTTS{frames: []voice.PCMFrame{mustFrame(t, []byte{1, 2})}},
		},
		AudioRuntime{
			Ingress: &fakeIngress{frames: []voice.PCMFrame{mustFrame(t, []byte{3, 4})}},
			Egress:  &fakeEgress{},
		},
	)
	if runner == nil {
		t.Fatal("runner = nil")
	}
}
