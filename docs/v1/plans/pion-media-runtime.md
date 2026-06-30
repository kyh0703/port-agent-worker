# Pion Media Runtime

## Goal

Pion tracks를 session `AudioRuntime`으로 변환하는 adapter factory를 만든다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-pion-media-runtime.md`

## Workspace

- Branch: feature/v1-pion-media-runtime
- Base: main
- Isolation: `.worktrees/feature-v1-pion-media-runtime`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: `pionrtp.FrameEncoder`를 public interface로 정리한다.
- Depends on:
  - none
- Write Scope:
  - `internal/adapters/pionrtp/**`
- Read Context:
  - `internal/adapters/pionrtp/egress.go`
- Checks:
  - `go test ./internal/adapters/pionrtp`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: Pion track config를 session `AudioRuntime`으로 변환하는 `media/pion` factory를 추가한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/media/pion/**`
- Read Context:
  - `internal/adapters/pionrtp/**`
  - `internal/application/session/assembler.go`
- Checks:
  - `go test ./internal/adapters/media/pion`
- Parallel-safe: no

### Task T3

- [x] Complete
- Goal: 전체 검증을 완료한다.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `internal/adapters/pionrtp/**`
  - `internal/adapters/media/pion/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no
