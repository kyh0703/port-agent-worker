# Turn Aware Session

## Goal

PCM stream을 STT와 VAD로 fan-out하고, VAD activity event를 turn decision으로 세션 경계에서 관찰할 수 있게 한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-turn-aware-session.md`
- `docs/v1/designs/2026-06-30-v1-voice-activity-flow.md`

## Workspace

- Branch: feature/v1-turn-aware-session
- Base: main
- Isolation: `.worktrees/feature-v1-turn-aware-session`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: `internal/application/session`에 turn runtime type과 assembler를 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/application/session/**`
- Read Context:
  - `internal/application/session/assembler.go`
  - `internal/application/ports/vad.go`
  - `internal/application/turn/**`
- Checks:
  - `go test ./internal/application/session`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: turn-aware orchestrator가 PCM fan-out, VAD event 처리, decision handler 호출을 수행하게 한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/application/session/**`
- Read Context:
  - `internal/application/session/orchestrator.go`
  - `internal/application/turn/**`
- Checks:
  - `go test ./internal/application/session`
- Parallel-safe: no

### Task T3

- [x] Complete
- Goal: 전체 Go 검증과 build 검증을 완료한다.
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
