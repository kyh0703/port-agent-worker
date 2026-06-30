---
feature: session-runner
created_at: 2026-06-30T14:05:00+09:00
---

# Session Runner

## Goal

`RunTurn` 단위 orchestrator를 반복 실행 가능한 session runner로 감싼다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-worker-runtime-wiring.md`
- Existing system facts:
  - `session.Orchestrator.RunTurn(ctx)`는 한 번의 STT -> LLM -> TTS -> egress 흐름만 실행한다.
  - `cmd/worker`는 provider wiring만 수행하고 실제 session loop는 없다.
  - Pion track acquisition은 아직 없다.
- User brief:
  - STT/LLM/TTS/Pion worker를 순차적으로 붙인다.
  - DDD + Hexagonal Architecture를 유지한다.

## Plan Handoff

### Scope for Planning

- `internal/application/session`에 반복 실행 runner를 추가한다.
- runner는 context cancel을 존중한다.
- `ErrNoFinalTranscript`는 idle/no turn으로 취급하고 다음 turn을 기다릴 수 있게 한다.
- fatal error는 caller에게 반환한다.
- fake orchestrator 기반 테스트로 반복 실행, cancel, fatal error를 검증한다.

### Success Criteria

- session runner가 context cancel 전까지 turn executor를 반복 호출한다.
- `ErrNoFinalTranscript`는 fatal error로 처리하지 않는다.
- fatal error는 wrapping 없이 caller가 식별 가능하게 반환된다.
- domain/application 계층에 adapter 타입이 노출되지 않는다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- Pion track acquisition.
- provider runtime integration.
- barge-in/endpointing.
- concurrency worker pool.

### Open Questions

- runner backoff 값은 MVP에서 짧은 고정 delay로 시작하고 운영 지표 후 조정한다.

### Suggested Validation

- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/application/session/**`
- Shared files to avoid touching in parallel:
  - `internal/application/session/orchestrator.go`
- Likely sequential dependencies:
  - runner interface -> loop implementation -> tests.

