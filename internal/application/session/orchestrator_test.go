package session

import (
	"context"
	"testing"
	"time"

	"port-voice-pipeline/internal/domain/voice"
)

func TestOrchestratorRunTurnWritesSynthesizedAudio(t *testing.T) {
	ctx := context.Background()
	frame := mustFrame(t, []byte{1, 2, 3, 4})

	ingress := &fakeIngress{frames: []voice.PCMFrame{frame}}
	stt := &fakeSTT{transcripts: []voice.Transcript{
		{Text: "partial", IsFinal: false},
		{Text: "hello", IsFinal: true},
	}}
	llm := &fakeLLM{response: voice.AssistantResponse{Text: "hi there"}}
	ttsFrame := mustFrame(t, []byte{5, 6, 7, 8})
	tts := &fakeTTS{frames: []voice.PCMFrame{ttsFrame}}
	egress := &fakeEgress{}

	orchestrator := NewOrchestrator(ingress, egress, stt, llm, tts)

	if err := orchestrator.RunTurn(ctx); err != nil {
		t.Fatalf("RunTurn() error = %v", err)
	}

	if llm.received.Text != "hello" {
		t.Fatalf("llm received text = %q, want %q", llm.received.Text, "hello")
	}

	if !egress.flushed {
		t.Fatal("egress was not flushed")
	}

	if len(egress.frames) != 1 {
		t.Fatalf("egress wrote %d frames, want 1", len(egress.frames))
	}

	if string(egress.frames[0].Data) != string(ttsFrame.Data) {
		t.Fatalf("egress frame data = %v, want %v", egress.frames[0].Data, ttsFrame.Data)
	}
}

func mustFrame(t *testing.T, data []byte) voice.PCMFrame {
	t.Helper()

	frame, err := voice.NewPCMFrame(data, 16000, 1, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("NewPCMFrame() error = %v", err)
	}
	return frame
}

type fakeIngress struct {
	frames []voice.PCMFrame
}

func (f *fakeIngress) PCMFrames(context.Context) (<-chan voice.PCMFrame, error) {
	out := make(chan voice.PCMFrame, len(f.frames))
	for _, frame := range f.frames {
		out <- frame
	}
	close(out)
	return out, nil
}

type fakeSTT struct {
	transcripts []voice.Transcript
}

func (f *fakeSTT) Transcribe(context.Context, <-chan voice.PCMFrame) (<-chan voice.Transcript, error) {
	out := make(chan voice.Transcript, len(f.transcripts))
	for _, transcript := range f.transcripts {
		out <- transcript
	}
	close(out)
	return out, nil
}

type fakeLLM struct {
	response voice.AssistantResponse
	received voice.UserUtterance
}

func (f *fakeLLM) Generate(_ context.Context, utterance voice.UserUtterance) (voice.AssistantResponse, error) {
	f.received = utterance
	return f.response, nil
}

type fakeTTS struct {
	frames []voice.PCMFrame
}

func (f *fakeTTS) Synthesize(context.Context, voice.AssistantResponse) (<-chan voice.PCMFrame, error) {
	out := make(chan voice.PCMFrame, len(f.frames))
	for _, frame := range f.frames {
		out <- frame
	}
	close(out)
	return out, nil
}

type fakeEgress struct {
	frames  []voice.PCMFrame
	flushed bool
}

func (f *fakeEgress) WritePCM(_ context.Context, frame voice.PCMFrame) error {
	f.frames = append(f.frames, frame)
	return nil
}

func (f *fakeEgress) Flush(context.Context) error {
	f.flushed = true
	return nil
}

func (f *fakeEgress) Close() error {
	return nil
}
