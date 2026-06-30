# Turn Controller

## Goal

Bot/user speech 이벤트 기반 barge-in과 endpointing 판단을 application 계층 순수 상태기계로 구현한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-turn-controller.md`

## Workspace

- Branch: feature/v1-turn-controller
- Base: main
- Isolation: `.worktrees/feature-v1-turn-controller`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [ ] Complete
- Goal: `internal/application/turn`에 barge-in/endpointing controller와 테스트를 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/application/turn/**`
- Read Context:
  - `docs/ARCHITECTURE.md`
  - `internal/application/session/**`
- Checks:
  - `go test ./internal/application/turn`
- Parallel-safe: no

### Task T2

- [ ] Complete
- Goal: 전체 Go 검증과 build 검증을 완료한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/application/turn/**`
- Read Context:
  - `internal/application/turn/**`
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no

