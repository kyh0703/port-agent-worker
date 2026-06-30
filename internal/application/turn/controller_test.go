package turn

import (
	"context"
	"testing"
	"time"
)

func TestControllerBargeInWhenUserStartsDuringBotSpeech(t *testing.T) {
	controller := NewController(Config{StopDelay: time.Second}, nil)

	controller.BotStarted()
	decision := controller.UserSpeechStarted(time.Unix(10, 0))

	if !decision.BargeIn {
		t.Fatal("expected barge-in when user starts during bot speech")
	}
	if controller.BotSpeaking() {
		t.Fatal("expected bot speech to stop after barge-in")
	}
}

func TestControllerEndpointsAfterBotStopsAndSilencePasses(t *testing.T) {
	controller := NewController(Config{StopDelay: time.Second}, nil)
	start := time.Unix(10, 0)

	controller.BotStarted()
	controller.BotStopped(start)

	got, err := controller.ShouldEndpoint(context.Background(), start.Add(999*time.Millisecond))
	if err != nil {
		t.Fatalf("ShouldEndpoint returned error: %v", err)
	}
	if got {
		t.Fatal("expected no endpoint before stop delay")
	}

	got, err = controller.ShouldEndpoint(context.Background(), start.Add(time.Second))
	if err != nil {
		t.Fatalf("ShouldEndpoint returned error: %v", err)
	}
	if !got {
		t.Fatal("expected endpoint after stop delay")
	}
}

func TestControllerDoesNotEndpointWhileUserSpeaking(t *testing.T) {
	controller := NewController(Config{StopDelay: time.Second}, nil)
	start := time.Unix(10, 0)

	controller.BotStopped(start)
	controller.UserSpeechStarted(start.Add(100 * time.Millisecond))

	got, err := controller.ShouldEndpoint(context.Background(), start.Add(2*time.Second))
	if err != nil {
		t.Fatalf("ShouldEndpoint returned error: %v", err)
	}
	if got {
		t.Fatal("expected no endpoint while user is speaking")
	}
}

func TestControllerRestartsSilenceAfterUserStops(t *testing.T) {
	controller := NewController(Config{StopDelay: time.Second}, nil)
	start := time.Unix(10, 0)

	controller.BotStopped(start)
	controller.UserSpeechStarted(start.Add(200 * time.Millisecond))
	controller.UserSpeechStopped(start.Add(500 * time.Millisecond))

	got, err := controller.ShouldEndpoint(context.Background(), start.Add(1200*time.Millisecond))
	if err != nil {
		t.Fatalf("ShouldEndpoint returned error: %v", err)
	}
	if got {
		t.Fatal("expected no endpoint before restarted silence reaches stop delay")
	}

	got, err = controller.ShouldEndpoint(context.Background(), start.Add(1500*time.Millisecond))
	if err != nil {
		t.Fatalf("ShouldEndpoint returned error: %v", err)
	}
	if !got {
		t.Fatal("expected endpoint after restarted silence reaches stop delay")
	}
}

func TestControllerUsesDefaultStopDelay(t *testing.T) {
	controller := NewController(Config{}, nil)
	start := time.Unix(10, 0)

	controller.BotStopped(start)

	got, err := controller.ShouldEndpoint(context.Background(), start.Add(DefaultStopDelay))
	if err != nil {
		t.Fatalf("ShouldEndpoint returned error: %v", err)
	}
	if !got {
		t.Fatal("expected default stop delay to be used")
	}
}

func TestControllerUsesSmartTurnAnalyzerWhenEnabled(t *testing.T) {
	analyzer := &fixedAnalyzer{endpoint: true}
	controller := NewController(Config{
		StopDelay:        time.Second,
		SmartTurnEnabled: true,
	}, analyzer)
	start := time.Unix(10, 0)

	controller.BotStopped(start)

	got, err := controller.ShouldEndpoint(context.Background(), start.Add(100*time.Millisecond))
	if err != nil {
		t.Fatalf("ShouldEndpoint returned error: %v", err)
	}
	if !got {
		t.Fatal("expected smart turn analyzer to end the turn")
	}
	if analyzer.calls != 1 {
		t.Fatalf("expected analyzer to be called once, got %d", analyzer.calls)
	}
	if analyzer.last.SilenceDuration != 100*time.Millisecond {
		t.Fatalf("expected analyzer silence duration 100ms, got %s", analyzer.last.SilenceDuration)
	}
	if !analyzer.last.Endpointing {
		t.Fatal("expected analyzer snapshot to be endpointing")
	}
}

type fixedAnalyzer struct {
	endpoint bool
	calls    int
	last     Snapshot
}

func (a *fixedAnalyzer) ShouldEndpoint(_ context.Context, snapshot Snapshot) (bool, error) {
	a.calls++
	a.last = snapshot
	return a.endpoint, nil
}
