# Deepgram STT Adapter

## Goal

Deepgram Nova-3 streaming STT adapter를 추가해서 PCM frame stream을 transcript stream으로 변환한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-deepgram-stt-adapter.md`
- Deepgram streaming STT docs.

## Workspace

- Branch: feature/v1-deepgram-stt-adapter
- Base: main
- Isolation: `.worktrees/feature-v1-deepgram-stt-adapter`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: Deepgram adapter configuration and WebSocket URL construction을 구현한다.
- Depends on:
  - none
- Write Scope:
  - `internal/adapters/providers/deepgram/**`
- Read Context:
  - `internal/application/ports/providers.go`
  - `internal/domain/voice/**`
  - `docs/v1/designs/2026-06-30-v1-deepgram-stt-adapter.md`
- Checks:
  - `go test ./...`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: Deepgram response parser를 구현하고 final/interim transcript 테스트를 추가한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/providers/deepgram/**`
- Read Context:
  - `internal/domain/voice/conversation.go`
- Checks:
  - `go test ./internal/adapters/providers/deepgram`
- Parallel-safe: no

### Task T3

- [x] Complete
- Goal: `SpeechToText` 포트 구현체로 WebSocket audio send/receive loop를 연결한다.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `internal/adapters/providers/deepgram/**`
- Read Context:
  - `internal/application/ports/providers.go`
  - `internal/domain/voice/audio.go`
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no
