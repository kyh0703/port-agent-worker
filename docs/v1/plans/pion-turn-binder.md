# Pion Turn Binder

## Goal

Pion media runtime 조립 경로에서 turn-aware session runner를 선택할 수 있도록 binder를 확장한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-pion-turn-binder.md`
- `docs/v1/designs/2026-06-30-v1-pion-session-binder.md`
- `docs/v1/designs/2026-06-30-v1-turn-runtime-wiring.md`

## Workspace

- Branch: feature/v1-pion-turn-binder
- Base: main
- Isolation: `.worktrees/feature-v1-pion-turn-binder`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [ ] Complete
- Goal: `internal/adapters/media/pion`에 turn-aware runner binder와 테스트를 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/adapters/media/pion/**`
- Read Context:
  - `internal/adapters/media/pion/binder.go`
  - `internal/adapters/media/pion/runtime.go`
  - `internal/application/session/**`
- Checks:
  - `go test ./internal/adapters/media/pion`
- Parallel-safe: no

### Task T2

- [ ] Complete
- Goal: 전체 Go 검증과 build 검증을 완료한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/media/pion/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no

