package cartesia

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"port-voice-pipeline/internal/application/ports"
	"port-voice-pipeline/internal/domain/voice"
)

const defaultBaseURL = "https://api.cartesia.ai/tts/bytes"

var (
	ErrMissingAPIKey = errors.New("cartesia api key is required")
	ErrMissingVoice  = errors.New("cartesia voice id is required")
	ErrEmptyText     = errors.New("tts text is empty")
)

type Config struct {
	APIKey     string
	BaseURL    string
	Version    string
	ModelID    string
	VoiceID    string
	Language   string
	Encoding   string
	SampleRate int
	Channels   int
	ChunkSize  int
	HTTPClient *http.Client
}

func (c Config) withDefaults() Config {
	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}
	if c.Version == "" {
		c.Version = "2026-03-01"
	}
	if c.ModelID == "" {
		c.ModelID = "sonic-3.5"
	}
	if c.Language == "" {
		c.Language = "ko"
	}
	if c.Encoding == "" {
		c.Encoding = "pcm_s16le"
	}
	if c.SampleRate == 0 {
		c.SampleRate = 16000
	}
	if c.Channels == 0 {
		c.Channels = 1
	}
	if c.ChunkSize == 0 {
		c.ChunkSize = 3200
	}
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
	return c
}

type Client struct {
	config Config
}

var _ ports.TextToSpeech = (*Client)(nil)

func New(config Config) *Client {
	return &Client{config: config.withDefaults()}
}

func (c *Client) Synthesize(ctx context.Context, response voice.AssistantResponse) (<-chan voice.PCMFrame, error) {
	cfg := c.config.withDefaults()
	if cfg.APIKey == "" {
		return nil, ErrMissingAPIKey
	}
	if cfg.VoiceID == "" {
		return nil, ErrMissingVoice
	}

	text := strings.TrimSpace(response.Text)
	if text == "" {
		return nil, ErrEmptyText
	}

	body, err := json.Marshal(newRequest(cfg, text))
	if err != nil {
		return nil, fmt.Errorf("encode cartesia request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.BaseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create cartesia request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Cartesia-Version", cfg.Version)
	req.Header.Set("Content-Type", "application/json")

	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call cartesia: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("cartesia status %d: %s", resp.StatusCode, strings.TrimSpace(string(payload)))
	}

	frames := make(chan voice.PCMFrame)
	go c.streamFrames(resp.Body, frames)
	return frames, nil
}

func (c *Client) streamFrames(body io.ReadCloser, frames chan<- voice.PCMFrame) {
	defer close(frames)
	defer body.Close()

	cfg := c.config.withDefaults()
	buf := make([]byte, cfg.ChunkSize)
	for {
		n, err := body.Read(buf)
		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			frame, frameErr := voice.NewPCMFrame(data, cfg.SampleRate, cfg.Channels, frameDuration(n, cfg.SampleRate, cfg.Channels))
			if frameErr == nil {
				frames <- frame
			}
		}
		if err != nil {
			return
		}
	}
}

func frameDuration(byteCount int, sampleRate int, channels int) time.Duration {
	const bytesPerSample = 2
	samples := byteCount / bytesPerSample / channels
	if samples <= 0 {
		return time.Nanosecond
	}
	return time.Duration(samples) * time.Second / time.Duration(sampleRate)
}

type requestBody struct {
	ModelID    string       `json:"model_id"`
	Transcript string       `json:"transcript"`
	Voice      voiceRequest `json:"voice"`
	Output     outputFormat `json:"output_format"`
	Language   string       `json:"language,omitempty"`
}

type voiceRequest struct {
	Mode string `json:"mode"`
	ID   string `json:"id"`
}

type outputFormat struct {
	Container  string `json:"container"`
	Encoding   string `json:"encoding,omitempty"`
	SampleRate int    `json:"sample_rate,omitempty"`
}

func newRequest(cfg Config, text string) requestBody {
	return requestBody{
		ModelID:    cfg.ModelID,
		Transcript: text,
		Voice: voiceRequest{
			Mode: "id",
			ID:   cfg.VoiceID,
		},
		Output: outputFormat{
			Container:  "raw",
			Encoding:   cfg.Encoding,
			SampleRate: cfg.SampleRate,
		},
		Language: cfg.Language,
	}
}
