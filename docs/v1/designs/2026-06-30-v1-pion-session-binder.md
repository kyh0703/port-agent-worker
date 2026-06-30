---
feature: pion-session-binder
created_at: 2026-06-30T15:25:00+09:00
---

# Pion Session Binder

## Goal

Pion media runtime과 provider runtime을 session runner로 묶는 adapter-level binder를 추가한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-pion-media-runtime.md`
  - `docs/v1/designs/2026-06-30-v1-worker-session-wiring.md`
- Existing system facts:
  - `internal/adapters/media/pion.NewRuntime`은 Pion tracks를 `session.AudioRuntime`으로 만든다.
  - `session.NewRunnerFromRuntime`은 provider/audio runtime을 runner로 만든다.
  - provider factory는 adapter 계층에서 STT/LLM/TTS 포트를 만든다.
- User brief:
  - Pion SFU를 유지한다.
  - DDD + Hexagonal Architecture를 유지한다.

## Plan Handoff

### Scope for Planning

- `internal/adapters/media/pion`에 provider runtime + Pion tracks + encoder를 runner로 조립하는 binder를 추가한다.
- binder는 `session.ProviderRuntime`을 받아 provider adapter 구체 타입을 모르게 한다.
- Pion track/encoder validation은 기존 media runtime factory를 재사용한다.

### Success Criteria

- `pion.NewRunner`가 Pion track config와 provider runtime을 받아 `*session.Runner`를 반환한다.
- Pion/provider 세부 타입은 adapter/application 경계를 넘지 않는다.
- track/encoder 누락 시 기존 에러가 유지된다.
- fake encoder/provider 기반 테스트가 있다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- SFU signaling.
- production Opus encoder.
- `cmd/worker` 실제 Pion session 실행.
- PeerConnection participant 구현.

### Open Questions

- SFU integration 방식은 다음 feature에서 확정한다.

### Suggested Validation

- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/media/pion/**`
- Shared files to avoid touching in parallel:
  - none expected.
- Likely sequential dependencies:
  - binder function -> tests.

