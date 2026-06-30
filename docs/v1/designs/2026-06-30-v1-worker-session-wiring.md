---
feature: worker-session-wiring
created_at: 2026-06-30T14:30:00+09:00
---

# Worker Session Wiring

## Goal

`cmd/worker`가 provider runtime에서 session orchestrator/runner까지 조립할 수 있게 한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-worker-runtime-wiring.md`
  - `docs/v1/designs/2026-06-30-v1-session-runner.md`
- Existing system facts:
  - provider factory는 STT/LLM/TTS 포트를 만든다.
  - session runner는 TurnExecutor를 반복 실행한다.
  - Pion track acquisition은 아직 없다.
  - production Opus encoder는 아직 없다.
- User brief:
  - STT/LLM/TTS/Pion worker를 순차적으로 붙인다.
  - DDD + Hexagonal Architecture를 유지한다.

## Plan Handoff

### Scope for Planning

- pending media adapter를 추가해 Pion track 미연결 상태를 명시적 에러로 표현한다.
- provider runtime + audio ports를 session runner로 조립하는 assembler를 추가한다.
- `RUN_SESSION=true`일 때 `cmd/worker`가 runner를 실행하게 한다.
- 기본값은 wiring validation 후 대기하여 현재 운영 동작을 유지한다.

### Success Criteria

- `RUN_SESSION` config가 로드된다.
- session assembler가 `ports.AudioIngress`, `ports.AudioEgress`, STT/LLM/TTS로 runner를 만든다.
- pending media adapter는 `ErrMediaNotConfigured`를 반환한다.
- `cmd/worker`는 `RUN_SESSION=true`일 때 runner를 실행한다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- 실제 Pion track acquisition.
- production Opus encoder.
- SFU signaling.
- endpointing/barge-in.

### Open Questions

- Pion SFU 연결 방식이 정해지면 pending media adapter를 실제 Pion media adapter factory로 교체한다.

### Suggested Validation

- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/media/pending/**`
  - `internal/application/session/**`
  - `cmd/worker/**`
  - `internal/config/**`
- Shared files to avoid touching in parallel:
  - `cmd/worker/main.go`
  - `internal/config/config.go`
- Likely sequential dependencies:
  - config -> pending media adapter -> session assembler -> cmd wiring.

