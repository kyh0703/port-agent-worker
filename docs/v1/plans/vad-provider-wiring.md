# VAD Provider Wiring

## Goal

`VAD_PROVIDER` 설정으로 noop 또는 Silero VAD adapter를 선택할 수 있게 한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-vad-provider-wiring.md`
- `docs/v1/designs/2026-06-30-v1-turn-runtime-wiring.md`
- `docs/v1/designs/2026-06-30-v1-silero-vad-adapter.md`

## Workspace

- Branch: feature/v1-vad-provider-wiring
- Base: main
- Isolation: `.worktrees/feature-v1-vad-provider-wiring`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: config에 VAD provider와 Silero VAD 설정을 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/config/**`
- Read Context:
  - `internal/config/config.go`
  - `internal/config/config_test.go`
- Checks:
  - `go test ./internal/config`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: turn runtime factory가 noop/silero provider 선택과 error handling을 수행하게 한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/turn/**`
- Read Context:
  - `internal/adapters/turn/runtime.go`
  - `internal/adapters/vad/noop/**`
  - `internal/adapters/vad/silero/**`
- Checks:
  - `go test ./internal/adapters/turn`
- Parallel-safe: no

### Task T3

- [x] Complete
- Goal: worker entrypoint가 turn runtime wiring error를 처리하게 한다.
- Depends on:
  - T2
- Write Scope:
  - `cmd/worker/**`
- Read Context:
  - `cmd/worker/main.go`
  - `internal/adapters/turn/**`
- Checks:
  - `go test ./cmd/worker`
- Parallel-safe: no

### Task T4

- [x] Complete
- Goal: 전체 Go 검증과 build 검증을 완료한다.
- Depends on:
  - T1
  - T2
  - T3
- Write Scope:
  - `internal/config/**`
  - `internal/adapters/turn/**`
  - `cmd/worker/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no
