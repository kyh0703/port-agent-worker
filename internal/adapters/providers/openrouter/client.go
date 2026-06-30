package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"port-agent-worker/internal/application/ports"
	"port-agent-worker/internal/domain/voice"
)

const defaultBaseURL = "https://openrouter.ai/api/v1/chat/completions"

var (
	ErrMissingAPIKey  = errors.New("openrouter api key is required")
	ErrEmptyUtterance = errors.New("user utterance is empty")
	ErrEmptyResponse  = errors.New("openrouter response is empty")
)

type Config struct {
	APIKey       string
	BaseURL      string
	Model        string
	SystemPrompt string
	AppTitle     string
	SiteURL      string
	HTTPClient   *http.Client
}

func (c Config) withDefaults() Config {
	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}
	if c.Model == "" {
		c.Model = "google/gemini-2.5-flash-lite"
	}
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
	return c
}

type Client struct {
	config Config
}

var _ ports.LanguageModel = (*Client)(nil)

func New(config Config) *Client {
	return &Client{config: config.withDefaults()}
}

func (c *Client) Generate(ctx context.Context, utterance voice.UserUtterance) (voice.AssistantResponse, error) {
	cfg := c.config.withDefaults()
	if cfg.APIKey == "" {
		return voice.AssistantResponse{}, ErrMissingAPIKey
	}

	text := strings.TrimSpace(utterance.Text)
	if text == "" {
		return voice.AssistantResponse{}, ErrEmptyUtterance
	}

	payload, err := json.Marshal(newRequest(cfg, text))
	if err != nil {
		return voice.AssistantResponse{}, fmt.Errorf("encode openrouter request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.BaseURL, bytes.NewReader(payload))
	if err != nil {
		return voice.AssistantResponse{}, fmt.Errorf("create openrouter request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	if cfg.SiteURL != "" {
		req.Header.Set("HTTP-Referer", cfg.SiteURL)
	}
	if cfg.AppTitle != "" {
		req.Header.Set("X-Title", cfg.AppTitle)
	}

	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return voice.AssistantResponse{}, fmt.Errorf("call openrouter: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return voice.AssistantResponse{}, fmt.Errorf("read openrouter response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return voice.AssistantResponse{}, fmt.Errorf("openrouter status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	content, err := parseResponse(body)
	if err != nil {
		return voice.AssistantResponse{}, err
	}
	return voice.AssistantResponse{Text: content}, nil
}

type requestBody struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func newRequest(cfg Config, userText string) requestBody {
	messages := make([]message, 0, 2)
	if prompt := strings.TrimSpace(cfg.SystemPrompt); prompt != "" {
		messages = append(messages, message{Role: "system", Content: prompt})
	}
	messages = append(messages, message{Role: "user", Content: userText})
	return requestBody{
		Model:    cfg.Model,
		Messages: messages,
	}
}

type responseBody struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func parseResponse(payload []byte) (string, error) {
	var response responseBody
	if err := json.Unmarshal(payload, &response); err != nil {
		return "", fmt.Errorf("decode openrouter response: %w", err)
	}
	if len(response.Choices) == 0 {
		return "", ErrEmptyResponse
	}

	content := strings.TrimSpace(response.Choices[0].Message.Content)
	if content == "" {
		return "", ErrEmptyResponse
	}
	return content, nil
}
