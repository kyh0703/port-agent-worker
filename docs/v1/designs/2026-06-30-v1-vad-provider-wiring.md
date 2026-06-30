---
feature: vad-provider-wiring
created_at: 2026-06-30T12:56:45+09:00
---

# VAD Provider Wiring

## Goal

turn runtime factory에서 `VAD_PROVIDER` 설정을 해석해 noop 또는 Silero VAD adapter를 선택할 수 있게 한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-turn-runtime-wiring.md`
  - `docs/v1/designs/2026-06-30-v1-silero-vad-adapter.md`
- Existing system facts:
  - `internal/adapters/turn.NewRuntime`은 항상 noop VAD를 사용한다.
  - `internal/adapters/vad/silero.Detector`는 `silero.Engine`을 주입받아 동작한다.
  - 실제 ONNX engine 구현은 아직 없다.
- User brief:
  - smart turn off 상태에서는 VAD 묵음 기반 endpointing이 필요하다.
  - DDD + Hexagonal Architecture 경계를 유지한다.

## Plan Handoff

### Scope for Planning

- config에 VAD provider와 Silero threshold/min frame 설정을 추가한다.
- turn runtime factory가 `noop`과 `silero` provider를 선택하도록 확장한다.
- `silero` 선택 시 engine이 없으면 명시적 에러를 반환한다.
- worker entrypoint가 turn runtime wiring error를 처리한다.

### Success Criteria

- 기본 `VAD_PROVIDER`는 `noop`이다.
- `VAD_PROVIDER=noop`은 기존 noop VAD runtime을 만든다.
- `VAD_PROVIDER=silero`와 engine 주입 시 Silero detector runtime을 만든다.
- `VAD_PROVIDER=silero`인데 engine이 없으면 명시적 에러를 반환한다.
- 알 수 없는 VAD provider는 명시적 에러를 반환한다.
- worker는 turn runtime wiring error 발생 시 실패 로그를 남기고 종료한다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- 실제 ONNX runtime engine 구현.
- model path/env 처리.
- resampling implementation.
- Silero를 기본 VAD로 변경.

### Suggested Validation

- `go test ./internal/adapters/turn`
- `go test ./internal/config`
- `go test ./cmd/worker`
- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/turn/**`
  - `internal/config/**`
  - `cmd/worker/**`
- Shared files to avoid touching in parallel:
  - `internal/adapters/turn/runtime.go`
  - `internal/config/config.go`
  - `cmd/worker/main.go`
- Likely sequential dependencies:
  - config fields -> turn factory -> worker error handling.

