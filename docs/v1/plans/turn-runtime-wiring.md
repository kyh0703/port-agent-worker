# Turn Runtime Wiring

## Goal

worker entrypoint에서 turn-aware runner를 선택적으로 조립하고, Silero 전 단계로 noop VAD adapter를 추가한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-turn-runtime-wiring.md`
- `docs/v1/designs/2026-06-30-v1-turn-aware-session.md`

## Workspace

- Branch: feature/v1-turn-runtime-wiring
- Base: main
- Isolation: `.worktrees/feature-v1-turn-runtime-wiring`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: noop VAD adapter와 테스트를 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/adapters/vad/noop/**`
- Read Context:
  - `internal/application/ports/vad.go`
  - `internal/domain/voice/**`
- Checks:
  - `go test ./internal/adapters/vad/noop`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: turn runtime factory와 logging decision handler를 추가한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/turn/**`
- Read Context:
  - `internal/application/session/turn_aware.go`
  - `internal/application/turn/**`
- Checks:
  - `go test ./internal/adapters/turn`
- Parallel-safe: no

### Task T3

- [x] Complete
- Goal: config와 `cmd/worker`에서 turn-aware runner 선택 wiring을 추가한다.
- Depends on:
  - T2
- Write Scope:
  - `internal/config/**`
  - `cmd/worker/**`
- Read Context:
  - `internal/adapters/turn/**`
  - `internal/application/session/**`
- Checks:
  - `go test ./internal/config`
  - `go test ./...`
- Parallel-safe: no

### Task T4

- [x] Complete
- Goal: 전체 Go 검증과 build 검증을 완료한다.
- Depends on:
  - T1
  - T2
  - T3
- Write Scope:
  - `internal/adapters/vad/noop/**`
  - `internal/adapters/turn/**`
  - `internal/config/**`
  - `cmd/worker/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no
