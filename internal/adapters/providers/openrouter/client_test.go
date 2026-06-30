package openrouter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"port-agent-worker/internal/domain/voice"
)

func TestGenerateRequiresConfig(t *testing.T) {
	_, err := New(Config{}).Generate(context.Background(), voice.UserUtterance{Text: "hello"})
	if err != ErrMissingAPIKey {
		t.Fatalf("Generate() error = %v, want %v", err, ErrMissingAPIKey)
	}

	_, err = New(Config{APIKey: "key"}).Generate(context.Background(), voice.UserUtterance{Text: "   "})
	if err != ErrEmptyUtterance {
		t.Fatalf("Generate() error = %v, want %v", err, ErrEmptyUtterance)
	}
}

func TestGeneratePostsChatCompletionAndParsesResponse(t *testing.T) {
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
			referer: r.Header.Get("HTTP-Referer"),
			title:   r.Header.Get("X-Title"),
			body:    body,
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":" hi there "}}]}`))
	}))
	defer server.Close()

	client := New(Config{
		APIKey:       "test-key",
		BaseURL:      server.URL,
		SystemPrompt: "be concise",
		AppTitle:     "port-agent-worker",
		SiteURL:      "https://example.test",
		HTTPClient:   server.Client(),
	})

	response, err := client.Generate(context.Background(), voice.UserUtterance{Text: " hello "})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if response.Text != "hi there" {
		t.Fatalf("response.Text = %q, want hi there", response.Text)
	}

	request := <-requests
	if request.auth != "Bearer test-key" {
		t.Fatalf("Authorization = %q, want Bearer test-key", request.auth)
	}
	if request.referer != "https://example.test" {
		t.Fatalf("HTTP-Referer = %q", request.referer)
	}
	if request.title != "port-agent-worker" {
		t.Fatalf("X-Title = %q", request.title)
	}
	if request.body.Model != "google/gemini-2.5-flash-lite" {
		t.Fatalf("Model = %q, want google/gemini-2.5-flash-lite", request.body.Model)
	}
	if len(request.body.Messages) != 2 {
		t.Fatalf("message count = %d, want 2", len(request.body.Messages))
	}
	if request.body.Messages[0].Role != "system" || request.body.Messages[0].Content != "be concise" {
		t.Fatalf("system message = %+v", request.body.Messages[0])
	}
	if request.body.Messages[1].Role != "user" || request.body.Messages[1].Content != "hello" {
		t.Fatalf("user message = %+v", request.body.Messages[1])
	}
}

func TestGenerateReturnsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	client := New(Config{
		APIKey:     "key",
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	})

	_, err := client.Generate(context.Background(), voice.UserUtterance{Text: "hello"})
	if err == nil {
		t.Fatal("Generate() error = nil, want error")
	}
}

func TestParseResponseRequiresContent(t *testing.T) {
	_, err := parseResponse([]byte(`{"choices":[{"message":{"role":"assistant","content":" "}}]}`))
	if err != ErrEmptyResponse {
		t.Fatalf("parseResponse() error = %v, want %v", err, ErrEmptyResponse)
	}
}

type capturedRequest struct {
	auth    string
	referer string
	title   string
	body    requestBody
}
