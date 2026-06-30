# Session Runner

## Goal

한 turn orchestration을 반복 실행 가능한 session runner로 감싸 session runtime의 application boundary를 만든다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-session-runner.md`

## Workspace

- Branch: feature/v1-session-runner
- Base: main
- Isolation: `.worktrees/feature-v1-session-runner`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: TurnExecutor interface와 Runner type을 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/application/session/**`
- Read Context:
  - `internal/application/session/orchestrator.go`
- Checks:
  - `go test ./internal/application/session`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: cancel, no-final-transcript, fatal error 동작 테스트를 추가한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/application/session/**`
- Read Context:
  - `internal/application/session/orchestrator_test.go`
- Checks:
  - `go test ./internal/application/session`
- Parallel-safe: no

### Task T3

- [x] Complete
- Goal: 전체 검증을 완료한다.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `internal/application/session/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no
