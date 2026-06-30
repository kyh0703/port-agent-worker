# Pion Session Binder

## Goal

Pion tracks와 provider runtime을 session runner로 조립하는 adapter-level binder를 만든다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-pion-session-binder.md`

## Workspace

- Branch: feature/v1-pion-session-binder
- Base: main
- Isolation: `.worktrees/feature-v1-pion-session-binder`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: `media/pion`에 provider runtime + Pion track config를 runner로 조립하는 binder를 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/adapters/media/pion/**`
- Read Context:
  - `internal/adapters/media/pion/runtime.go`
  - `internal/application/session/assembler.go`
- Checks:
  - `go test ./internal/adapters/media/pion`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: fake provider/encoder 기반 테스트와 전체 검증을 완료한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/media/pion/**`
- Read Context:
  - `internal/application/session/**`
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no
