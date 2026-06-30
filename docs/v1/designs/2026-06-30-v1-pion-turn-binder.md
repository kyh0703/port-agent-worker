---
feature: pion-turn-binder
created_at: 2026-06-30T12:46:31+09:00
---

# Pion Turn Binder

## Goal

Pion media runtime 조립 경로에서 turn-aware session runner를 선택할 수 있도록 adapter-level binder를 확장한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-pion-session-binder.md`
  - `docs/v1/designs/2026-06-30-v1-turn-runtime-wiring.md`
- Existing system facts:
  - `internal/adapters/media/pion.NewRunner`는 Pion track config와 provider runtime을 기본 session runner로 조립한다.
  - `session.NewTurnAwareRunnerFromRuntime`은 provider/audio/turn runtime을 turn-aware runner로 조립한다.
  - `cmd/worker` pending media 경로는 이미 `TURN_ENABLED`로 turn-aware runner를 선택할 수 있다.
- User brief:
  - Pion SFU를 유지한다.
  - worker 내부는 DDD + Hexagonal Architecture로 쪼갠다.

## Plan Handoff

### Scope for Planning

- `internal/adapters/media/pion`에 `NewTurnAwareRunner` binder를 추가한다.
- 기존 `NewRunner` 동작과 API는 유지한다.
- binder는 provider runtime, Pion track config, turn runtime을 받아 `*session.Runner`를 반환한다.
- Pion 타입은 adapter 계층 안에 머물고 application 계층으로 노출하지 않는다.

### Success Criteria

- `pion.NewTurnAwareRunner`가 Pion media runtime과 turn runtime을 조립한다.
- track/encoder 누락 시 기존 `NewRuntime` 에러가 유지된다.
- 기존 `pion.NewRunner` 테스트는 유지된다.
- fake provider/encoder/turn runtime 기반 테스트가 있다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- 실제 SFU signaling.
- production Opus encoder.
- Silero VAD adapter 구현.
- TTS/egress interrupt 구현.

### Suggested Validation

- `go test ./internal/adapters/media/pion`
- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/media/pion/**`
- Shared files to avoid touching in parallel:
  - `internal/adapters/media/pion/binder.go`
- Likely sequential dependencies:
  - binder function -> tests.

