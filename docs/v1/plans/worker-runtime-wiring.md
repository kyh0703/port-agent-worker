# Worker Runtime Wiring

## Goal

환경 설정을 provider adapter와 session orchestrator 조립에 연결해 worker runtime wiring을 준비한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-worker-runtime-wiring.md`

## Workspace

- Branch: feature/v1-worker-runtime-wiring
- Base: main
- Isolation: `.worktrees/feature-v1-worker-runtime-wiring`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [ ] Complete
- Goal: provider runtime config와 env loading 테스트를 추가한다.
- Depends on:
  - none
- Write Scope:
  - `internal/config/**`
- Read Context:
  - `internal/adapters/providers/*/client.go`
- Checks:
  - `go test ./internal/config`
- Parallel-safe: no

### Task T2

- [ ] Complete
- Goal: config를 STT/LLM/TTS port 구현체로 변환하는 provider factory를 추가한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/providers/**`
- Read Context:
  - `internal/application/ports/providers.go`
  - `internal/config/**`
- Checks:
  - `go test ./internal/adapters/providers`
- Parallel-safe: no

### Task T3

- [ ] Complete
- Goal: `cmd/worker`에서 provider wiring을 실행하고 실패를 exit code로 처리한다.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `cmd/worker/**`
- Read Context:
  - `internal/config/**`
  - `internal/adapters/providers/**`
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no

