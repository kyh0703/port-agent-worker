# Worker Session Wiring

## Goal

worker runtime provider와 audio port를 session runner까지 연결해 실행 경계를 만든다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-worker-session-wiring.md`

## Workspace

- Branch: feature/v1-worker-session-wiring
- Base: main
- Isolation: `.worktrees/feature-v1-worker-session-wiring`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: `RUN_SESSION` config와 pending media adapter를 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/config/**`
  - `internal/adapters/media/pending/**`
- Read Context:
  - `internal/application/ports/audio.go`
- Checks:
  - `go test ./internal/config`
  - `go test ./internal/adapters/media/pending`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: provider runtime과 media ports를 session runner로 조립하는 assembler를 추가한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/application/session/**`
- Read Context:
  - `internal/adapters/providers/factory.go`
  - `internal/application/session/orchestrator.go`
  - `internal/application/session/runner.go`
- Checks:
  - `go test ./internal/application/session`
- Parallel-safe: no

### Task T3

- [x] Complete
- Goal: `cmd/worker`가 `RUN_SESSION=true`일 때 runner를 실행하게 연결한다.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `cmd/worker/**`
- Read Context:
  - `internal/config/**`
  - `internal/application/session/**`
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no
