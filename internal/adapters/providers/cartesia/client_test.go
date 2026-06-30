package cartesia

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"port-agent-worker/internal/domain/voice"
)

func TestSynthesizeRequiresConfig(t *testing.T) {
	_, err := New(Config{}).Synthesize(context.Background(), voice.AssistantResponse{Text: "hello"})
	if err != ErrMissingAPIKey {
		t.Fatalf("Synthesize() error = %v, want %v", err, ErrMissingAPIKey)
	}

	_, err = New(Config{APIKey: "key"}).Synthesize(context.Background(), voice.AssistantResponse{Text: "hello"})
	if err != ErrMissingVoice {
		t.Fatalf("Synthesize() error = %v, want %v", err, ErrMissingVoice)
	}
}

func TestSynthesizePostsRequestAndStreamsPCMFrames(t *testing.T) {
	requests := make(chan capturedRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		payload, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("ReadAll() error = %v", err)
			return
		}

		var body requestBody
		if err := json.Unmarshal(payload, &body); err != nil {
			t.Errorf("Unmarshal() error = %v", err)
			return
		}

		requests <- capturedRequest{
			auth:    r.Header.Get("Authorization"),
			version: r.Header.Get("Cartesia-Version"),
			body:    body,
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	}))
	defer server.Close()

	client := New(Config{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		VoiceID:    "voice-id",
		ChunkSize:  4,
		HTTPClient: server.Client(),
	})

	frames, err := client.Synthesize(context.Background(), voice.AssistantResponse{Text: " hello "})
	if err != nil {
		t.Fatalf("Synthesize() error = %v", err)
	}

	request := <-requests
	if request.auth != "Bearer test-key" {
		t.Fatalf("Authorization = %q, want Bearer test-key", request.auth)
	}
	if request.version != "2026-03-01" {
		t.Fatalf("Cartesia-Version = %q, want 2026-03-01", request.version)
	}
	if request.body.ModelID != "sonic-3.5" {
		t.Fatalf("ModelID = %q, want sonic-3.5", request.body.ModelID)
	}
	if request.body.Transcript != "hello" {
		t.Fatalf("Transcript = %q, want hello", request.body.Transcript)
	}
	if request.body.Voice.Mode != "id" || request.body.Voice.ID != "voice-id" {
		t.Fatalf("Voice = %+v, want mode id with voice-id", request.body.Voice)
	}
	if request.body.Output.Container != "raw" || request.body.Output.Encoding != "pcm_s16le" || request.body.Output.SampleRate != 16000 {
		t.Fatalf("Output = %+v, want raw pcm_s16le 16000", request.body.Output)
	}
	if request.body.Language != "ko" {
		t.Fatalf("Language = %q, want ko", request.body.Language)
	}

	got := collectFrames(frames)
	if len(got) != 2 {
		t.Fatalf("frame count = %d, want 2", len(got))
	}
	if string(got[0].Data) != string([]byte{1, 2, 3, 4}) {
		t.Fatalf("first frame = %v, want [1 2 3 4]", got[0].Data)
	}
	if got[0].SampleRate != 16000 || got[0].Channels != 1 || got[0].Encoding != voice.PCMEncodingS16LE {
		t.Fatalf("first frame metadata = %+v", got[0])
	}
}

func TestSynthesizeReturnsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	client := New(Config{
		APIKey:     "key",
		BaseURL:    server.URL,
		VoiceID:    "voice-id",
		HTTPClient: server.Client(),
	})

	_, err := client.Synthesize(context.Background(), voice.AssistantResponse{Text: "hello"})
	if err == nil {
		t.Fatal("Synthesize() error = nil, want error")
	}
}

type capturedRequest struct {
	auth    string
	version string
	body    requestBody
}

func collectFrames(frames <-chan voice.PCMFrame) []voice.PCMFrame {
	var out []voice.PCMFrame
	for frame := range frames {
		out = append(out, frame)
	}
	return out
}
