# Pipeline Stability and Hexagonal Wiring Cleanup

## Goal

Fix turn-aware audio fan-out blocking behavior and clean provider wiring so adapter packages no longer depend on global config.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v2/designs/2026-06-30-v2-pipeline-stability-hex-wiring.md`

## Workspace

- Branch: fix/v2-pipeline-stability-hex-wiring
- Base: main
- Isolation: `.worktrees/fix-v2-pipeline-stability-hex-wiring`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: Reproduce and fix turn-aware fan-out so a stopped VAD consumer cannot block STT input completion.
- Depends on:
  - none
- Write Scope:
  - `internal/application/session/**`
- Read Context:
  - `internal/application/session/turn_aware.go`
  - `internal/application/session/turn_aware_test.go`
  - `docs/v2/designs/2026-06-30-v2-pipeline-stability-hex-wiring.md`
- Checks:
  - `go test ./internal/application/session`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: Refactor provider runtime construction so `internal/adapters/providers` does not import `internal/config`.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/providers/**`
  - `cmd/worker/**`
- Read Context:
  - `internal/adapters/providers/factory.go`
  - `internal/adapters/providers/factory_test.go`
  - `cmd/worker/main.go`
  - `cmd/worker/main_test.go`
  - `internal/config/config.go`
  - `docs/v2/designs/2026-06-30-v2-pipeline-stability-hex-wiring.md`
- Checks:
  - `go test ./internal/adapters/providers ./cmd/worker`
  - `go list -f '{{.ImportPath}} -> {{join .Imports " "}}' ./internal/adapters/providers`
- Parallel-safe: no

### Task T3

- [x] Complete
- Goal: Run full verification for the stability and wiring cleanup slice.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `internal/application/session/**`
  - `internal/adapters/providers/**`
  - `cmd/worker/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `go vet ./...`
  - `make build`
- Parallel-safe: no
