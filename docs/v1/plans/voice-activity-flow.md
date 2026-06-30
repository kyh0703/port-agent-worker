# Voice Activity Flow

## Goal

VAD speech activity event와 turn controller 연결 계약을 추가해 session/media 통합 전 단계를 만든다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-voice-activity-flow.md`
- `docs/v1/designs/2026-06-30-v1-turn-controller.md`

## Workspace

- Branch: feature/v1-voice-activity-flow
- Base: main
- Isolation: `.worktrees/feature-v1-voice-activity-flow`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: `internal/domain/voice`에 speech activity event value와 validation을 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/domain/voice/**`
- Read Context:
  - `internal/domain/voice/**`
- Checks:
  - `go test ./internal/domain/voice`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: `internal/application/ports`에 VAD adapter port를 추가한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/application/ports/**`
- Read Context:
  - `internal/domain/voice/**`
- Checks:
  - `go test ./internal/application/ports`
- Parallel-safe: no

### Task T3

- [x] Complete
- Goal: `internal/application/turn`에 speech activity event processor와 테스트를 추가한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/application/turn/**`
- Read Context:
  - `internal/application/turn/**`
  - `internal/domain/voice/**`
- Checks:
  - `go test ./internal/application/turn`
- Parallel-safe: no

### Task T4

- [x] Complete
- Goal: 전체 Go 검증과 build 검증을 완료한다.
- Depends on:
  - T1
  - T2
  - T3
- Write Scope:
  - `internal/domain/voice/**`
  - `internal/application/ports/**`
  - `internal/application/turn/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no
