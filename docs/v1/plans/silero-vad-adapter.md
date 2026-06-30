# Silero VAD Adapter

## Goal

Silero ONNX engine을 나중에 붙일 수 있는 VAD adapter skeleton과 speech/silence event 상태 전이를 구현한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-silero-vad-adapter.md`
- `docs/v1/designs/2026-06-30-v1-voice-activity-flow.md`
- `docs/v1/designs/2026-06-30-v1-turn-runtime-wiring.md`

## Workspace

- Branch: feature/v1-silero-vad-adapter
- Base: main
- Isolation: `.worktrees/feature-v1-silero-vad-adapter`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [ ] Complete
- Goal: `internal/adapters/vad/silero`에 engine 기반 detector와 테스트를 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/adapters/vad/silero/**`
- Read Context:
  - `internal/application/ports/vad.go`
  - `internal/domain/voice/**`
  - `internal/adapters/vad/noop/**`
- Checks:
  - `go test ./internal/adapters/vad/silero`
- Parallel-safe: no

### Task T2

- [ ] Complete
- Goal: 전체 Go 검증과 build 검증을 완료한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/vad/silero/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no

