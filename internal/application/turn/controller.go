package turn

import (
	"context"
	"time"
)

const DefaultStopDelay = 700 * time.Millisecond

type Config struct {
	StopDelay        time.Duration
	SmartTurnEnabled bool
}

type Decision struct {
	BargeIn  bool
	Endpoint bool
}

type Snapshot struct {
	BotSpeaking      bool
	UserSpeaking     bool
	Endpointing      bool
	SilenceStartedAt time.Time
	SilenceDuration  time.Duration
}

type SmartTurnAnalyzer interface {
	ShouldEndpoint(ctx context.Context, snapshot Snapshot) (bool, error)
}

type Controller struct {
	stopDelay        time.Duration
	smartTurnEnabled bool
	analyzer         SmartTurnAnalyzer

	botSpeaking    bool
	userSpeaking   bool
	endpointing    bool
	silenceStarted time.Time
}

func NewController(config Config, analyzer SmartTurnAnalyzer) *Controller {
	stopDelay := config.StopDelay
	if stopDelay <= 0 {
		stopDelay = DefaultStopDelay
	}

	return &Controller{
		stopDelay:        stopDelay,
		smartTurnEnabled: config.SmartTurnEnabled,
		analyzer:         analyzer,
	}
}

func (c *Controller) BotStarted() {
	c.botSpeaking = true
	c.endpointing = false
	c.silenceStarted = time.Time{}
}

func (c *Controller) BotStopped(now time.Time) {
	c.botSpeaking = false
	c.endpointing = true
	if !c.userSpeaking {
		c.silenceStarted = now
	}
}

func (c *Controller) UserSpeechStarted(time.Time) Decision {
	c.userSpeaking = true
	c.endpointing = false
	c.silenceStarted = time.Time{}

	if !c.botSpeaking {
		return Decision{}
	}

	c.botSpeaking = false
	return Decision{BargeIn: true}
}

func (c *Controller) UserSpeechStopped(now time.Time) {
	c.userSpeaking = false
	c.endpointing = true
	c.silenceStarted = now
}

func (c *Controller) ShouldEndpoint(ctx context.Context, now time.Time) (bool, error) {
	if !c.endpointing || c.userSpeaking || c.silenceStarted.IsZero() {
		return false, nil
	}

	snapshot := c.Snapshot(now)
	if c.smartTurnEnabled && c.analyzer != nil {
		endpoint, err := c.analyzer.ShouldEndpoint(ctx, snapshot)
		if err != nil {
			return false, err
		}
		if endpoint {
			return true, nil
		}
	}

	return snapshot.SilenceDuration >= c.stopDelay, nil
}

func (c *Controller) BotSpeaking() bool {
	return c.botSpeaking
}

func (c *Controller) Snapshot(now time.Time) Snapshot {
	silenceDuration := time.Duration(0)
	if !c.silenceStarted.IsZero() && now.After(c.silenceStarted) {
		silenceDuration = now.Sub(c.silenceStarted)
	}

	return Snapshot{
		BotSpeaking:      c.botSpeaking,
		UserSpeaking:     c.userSpeaking,
		Endpointing:      c.endpointing,
		SilenceStartedAt: c.silenceStarted,
		SilenceDuration:  silenceDuration,
	}
}
