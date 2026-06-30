---
feature: turn-runtime-wiring
created_at: 2026-06-30T12:39:59+09:00
---

# Turn Runtime Wiring

## Goal

worker entrypoint에서 turn-aware session runner를 선택적으로 조립할 수 있게 하고, 실제 Silero 구현 전까지 사용할 noop VAD adapter를 추가한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-turn-aware-session.md`
  - `docs/v1/designs/2026-06-30-v1-voice-activity-flow.md`
- Existing system facts:
  - `session.NewTurnAwareRunnerFromRuntime`은 `TurnRuntime`을 받아 VAD event를 turn decision으로 처리한다.
  - `cmd/worker`는 아직 항상 `session.NewRunnerFromRuntime`만 사용한다.
  - `SmartTurnEnabled` config는 있지만 실제 LocalSmartTurnAnalyzerV3는 MVP deferred이다.
- User brief:
  - 한 worker가 STT, LLM, TTS, 송출 흐름을 가져간다.
  - 내부는 DDD + Hexagonal Architecture로 쪼갠다.

## Plan Handoff

### Scope for Planning

- `internal/adapters/vad/noop`에 `ports.VoiceActivityDetector` 구현을 추가한다.
- `internal/adapters/turn`에 `session.TurnRuntime` factory와 logging decision handler를 추가한다.
- config에 turn runtime 활성화 플래그를 추가한다.
- `cmd/worker`에서 turn runtime 활성화 시 `NewTurnAwareRunnerFromRuntime`을 사용한다.

### Success Criteria

- noop VAD는 PCM frame stream을 drain하고 빈 speech activity stream을 닫는다.
- turn runtime factory는 noop VAD와 `turn.ActivityProcessor`를 조립한다.
- `TURN_ENABLED=true`일 때 worker가 turn-aware runner를 조립한다.
- 기본값에서는 기존 runner 경로가 유지된다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- Silero VAD adapter 구현.
- LocalSmartTurnAnalyzerV3 구현.
- TTS/egress interrupt 구현.
- 실제 Pion track acquire/signaling.

### Open Questions

- 실제 VAD provider 선택 env는 Silero adapter가 추가될 때 확정한다.

### Suggested Validation

- `go test ./internal/adapters/vad/noop`
- `go test ./internal/adapters/turn`
- `go test ./internal/config`
- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/vad/noop/**`
  - `internal/adapters/turn/**`
  - `internal/config/**`
  - `cmd/worker/**`
- Shared files to avoid touching in parallel:
  - `cmd/worker/main.go`
  - `internal/config/config.go`
- Likely sequential dependencies:
  - noop VAD -> turn runtime factory -> worker wiring.

