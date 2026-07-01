package deepgram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"port-voice-pipeline/internal/application/ports"
	"port-voice-pipeline/internal/domain/voice"
)

const defaultBaseURL = "wss://api.deepgram.com/v1/listen"

var (
	ErrMissingAPIKey       = errors.New("deepgram api key is required")
	ErrUnsupportedPCMFrame = errors.New("unsupported pcm frame")
)

type Config struct {
	APIKey         string
	BaseURL        string
	Model          string
	Language       string
	Encoding       string
	SampleRate     int
	Channels       int
	InterimResults bool
	SmartFormat    bool
}

func (c Config) withDefaults() Config {
	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}
	if c.Model == "" {
		c.Model = "nova-3"
	}
	if c.Language == "" {
		c.Language = "ko"
	}
	if c.Encoding == "" {
		c.Encoding = "linear16"
	}
	if c.SampleRate == 0 {
		c.SampleRate = 16000
	}
	if c.Channels == 0 {
		c.Channels = 1
	}
	return c
}

type Client struct {
	config Config
	dialer websocketDialer
}

var _ ports.SpeechToText = (*Client)(nil)

type websocketDialer interface {
	DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)
}

func New(config Config) *Client {
	return &Client{
		config: config.withDefaults(),
		dialer: websocket.DefaultDialer,
	}
}

func (c *Client) Transcribe(ctx context.Context, audio <-chan voice.PCMFrame) (<-chan voice.Transcript, error) {
	if c.config.APIKey == "" {
		return nil, ErrMissingAPIKey
	}

	conn, _, err := c.dialer.DialContext(ctx, c.listenURL(), c.headers())
	if err != nil {
		return nil, fmt.Errorf("connect deepgram: %w", err)
	}

	transcripts := make(chan voice.Transcript)
	go c.sendAudio(ctx, conn, audio)
	go c.receiveTranscripts(ctx, conn, transcripts)

	return transcripts, nil
}

func (c *Client) listenURL() string {
	cfg := c.config.withDefaults()
	endpoint, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return cfg.BaseURL
	}

	query := endpoint.Query()
	query.Set("model", cfg.Model)
	query.Set("language", cfg.Language)
	query.Set("encoding", cfg.Encoding)
	query.Set("sample_rate", strconv.Itoa(cfg.SampleRate))
	query.Set("channels", strconv.Itoa(cfg.Channels))
	query.Set("interim_results", strconv.FormatBool(cfg.InterimResults))
	query.Set("smart_format", strconv.FormatBool(cfg.SmartFormat))
	endpoint.RawQuery = query.Encode()
	return endpoint.String()
}

func (c *Client) headers() http.Header {
	headers := make(http.Header)
	headers.Set("Authorization", "Token "+c.config.APIKey)
	return headers
}

func (c *Client) sendAudio(ctx context.Context, conn *websocket.Conn, audio <-chan voice.PCMFrame) {
	defer conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))

	for {
		select {
		case <-ctx.Done():
			return
		case frame, ok := <-audio:
			if !ok {
				_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"CloseStream"}`))
				return
			}
			if !c.validFrame(frame) {
				_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ErrUnsupportedPCMFrame.Error()), time.Now().Add(time.Second))
				return
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, frame.Data); err != nil {
				return
			}
		}
	}
}

func (c *Client) validFrame(frame voice.PCMFrame) bool {
	cfg := c.config.withDefaults()
	return frame.Encoding == voice.PCMEncodingS16LE &&
		frame.SampleRate == cfg.SampleRate &&
		frame.Channels == cfg.Channels &&
		len(frame.Data) > 0
}

func (c *Client) receiveTranscripts(ctx context.Context, conn *websocket.Conn, transcripts chan<- voice.Transcript) {
	defer close(transcripts)
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, payload, err := conn.ReadMessage()
		if err != nil {
			return
		}

		transcript, ok := parseTranscript(payload)
		if !ok {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case transcripts <- transcript:
		}
	}
}

type resultMessage struct {
	Type    string `json:"type"`
	IsFinal bool   `json:"is_final"`
	Channel struct {
		Alternatives []struct {
			Transcript string `json:"transcript"`
		} `json:"alternatives"`
	} `json:"channel"`
}

func parseTranscript(payload []byte) (voice.Transcript, bool) {
	var message resultMessage
	if err := json.Unmarshal(payload, &message); err != nil {
		return voice.Transcript{}, false
	}
	if message.Type != "" && message.Type != "Results" {
		return voice.Transcript{}, false
	}
	if len(message.Channel.Alternatives) == 0 {
		return voice.Transcript{}, false
	}

	transcript := voice.Transcript{
		Text:    message.Channel.Alternatives[0].Transcript,
		IsFinal: message.IsFinal,
	}
	if text, ok := transcript.FinalText(); ok {
		transcript.Text = text
		return transcript, true
	}
	if transcript.Text == "" {
		return voice.Transcript{}, false
	}
	return transcript, true
}
