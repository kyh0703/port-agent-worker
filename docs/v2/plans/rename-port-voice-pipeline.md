# Rename Port Voice Pipeline

## Goal

Rename the Go module, runtime identity, and release image identity from `port-agent-worker` to `port-voice-pipeline` without changing runtime behavior.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v2/designs/2026-07-01-v2-rename-port-voice-pipeline.md`

## Workspace

- Branch: chore/v2-rename-port-voice-pipeline
- Base: main
- Isolation: `.worktrees/chore-v2-rename-port-voice-pipeline`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: Rename module imports, runtime/build/release identity, and docs from `port-agent-worker` to `port-voice-pipeline`.
- Depends on:
  - none
- Write Scope:
  - `go.mod`
  - `cmd/**`
  - `internal/**`
  - `.github/workflows/**`
  - `README.md`
  - `Makefile`
  - `Dockerfile`
  - `rails.toml`
  - `docs/**`
- Read Context:
  - `docs/v2/designs/2026-07-01-v2-rename-port-voice-pipeline.md`
  - `README.md`
  - `go.mod`
  - `Makefile`
  - `Dockerfile`
  - `rails.toml`
  - `.github/workflows/release-image.yml`
  - `cmd/worker/main.go`
- Checks:
  - `go test ./...`
  - `make build`
  - `! rg -n "port-agent-worker|voice agent worker|agent worker" README.md go.mod Makefile Dockerfile rails.toml .github cmd internal`
  - `! rg -n "port-agent-worker|voice agent worker|agent worker|worker-side" docs/ARCHITECTURE.md docs/STATE.md docs/ROADMAP.md docs/v2/designs/2026-06-30-v2-tag-release-image.md docs/v2/plans/tag-release-image.md`
- Parallel-safe: no

## Notes

- Keep `cmd/worker` as the entrypoint path for this slice.
