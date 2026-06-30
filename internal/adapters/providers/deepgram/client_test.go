package deepgram

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"port-agent-worker/internal/domain/voice"
)

func TestListenURLUsesNova3PCMDefaults(t *testing.T) {
	client := New(Config{
		APIKey:         "key",
		BaseURL:        "wss://example.test/v1/listen",
		InterimResults: true,
		SmartFormat:    true,
	})

	url := client.listenURL()
	for _, want := range []string{
		"model=nova-3",
		"language=ko",
		"encoding=linear16",
		"sample_rate=16000",
		"channels=1",
		"interim_results=true",
		"smart_format=true",
	} {
		if !strings.Contains(url, want) {
			t.Fatalf("listenURL() = %q, missing %q", url, want)
		}
	}
}

func TestTranscribeRequiresAPIKey(t *testing.T) {
	client := New(Config{})

	_, err := client.Transcribe(context.Background(), make(chan voice.PCMFrame))
	if err != ErrMissingAPIKey {
		t.Fatalf("Transcribe() error = %v, want %v", err, ErrMissingAPIKey)
	}
}

func TestParseTranscript(t *testing.T) {
	payload := []byte(`{"type":"Results","is_final":true,"channel":{"alternatives":[{"transcript":" hello "} ]}}`)

	transcript, ok := parseTranscript(payload)
	if !ok {
		t.Fatal("parseTranscript() ok = false")
	}
	if transcript.Text != "hello" {
		t.Fatalf("Text = %q, want hello", transcript.Text)
	}
	if !transcript.IsFinal {
		t.Fatal("IsFinal = false, want true")
	}
}

func TestParseInterimTranscript(t *testing.T) {
	payload := []byte(`{"type":"Results","is_final":false,"channel":{"alternatives":[{"transcript":"hel"}]}}`)

	transcript, ok := parseTranscript(payload)
	if !ok {
		t.Fatal("parseTranscript() ok = false")
	}
	if transcript.Text != "hel" {
		t.Fatalf("Text = %q, want hel", transcript.Text)
	}
	if transcript.IsFinal {
		t.Fatal("IsFinal = true, want false")
	}
}

func TestTranscribeStreamsAudioAndReceivesTranscript(t *testing.T) {
	upgrader := websocket.Upgrader{}
	audioReceived := make(chan []byte, 1)
	authReceived := make(chan string, 1)
	queryReceived := make(chan string, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authReceived <- r.Header.Get("Authorization")
		queryReceived <- r.URL.RawQuery

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("Upgrade() error = %v", err)
			return
		}
		defer conn.Close()

		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			t.Errorf("ReadMessage() error = %v", err)
			return
		}
		if messageType != websocket.BinaryMessage {
			t.Errorf("message type = %d, want binary", messageType)
			return
		}
		audioReceived <- payload

		if err := conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"Results","is_final":true,"channel":{"alternatives":[{"transcript":"hello"}]}}`)); err != nil {
			t.Errorf("WriteMessage() error = %v", err)
			return
		}
	}))
	defer server.Close()

	audio := make(chan voice.PCMFrame, 1)
	audio <- mustFrame(t, []byte{1, 2, 3, 4})
	close(audio)

	client := New(Config{
		APIKey:         "test-key",
		BaseURL:        "ws" + strings.TrimPrefix(server.URL, "http"),
		InterimResults: true,
		SmartFormat:    true,
	})

	transcripts, err := client.Transcribe(context.Background(), audio)
	if err != nil {
		t.Fatalf("Transcribe() error = %v", err)
	}

	select {
	case got := <-audioReceived:
		if string(got) != string([]byte{1, 2, 3, 4}) {
			t.Fatalf("audio payload = %v, want [1 2 3 4]", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for audio")
	}

	select {
	case got := <-authReceived:
		if got != "Token test-key" {
			t.Fatalf("Authorization = %q, want Token test-key", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for auth")
	}

	select {
	case got := <-queryReceived:
		for _, want := range []string{"model=nova-3", "encoding=linear16", "sample_rate=16000"} {
			if !strings.Contains(got, want) {
				t.Fatalf("query = %q, missing %q", got, want)
			}
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for query")
	}

	select {
	case transcript := <-transcripts:
		if transcript.Text != "hello" || !transcript.IsFinal {
			t.Fatalf("transcript = %+v, want final hello", transcript)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for transcript")
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
