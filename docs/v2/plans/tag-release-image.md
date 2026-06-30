# Tag Release Image

## Goal

Add a tag-triggered GitHub Actions workflow that builds the worker container image and pushes it to GHCR.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v2/designs/2026-06-30-v2-tag-release-image.md`

## Workspace

- Branch: feat/v2-tag-release-image
- Base: main
- Isolation: `.worktrees/feat-v2-tag-release-image`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: Add a production Dockerfile and docker ignore rules for the Go worker.
- Depends on:
  - none
- Write Scope:
  - `Dockerfile`
  - `.dockerignore`
- Read Context:
  - `go.mod`
  - `Makefile`
  - `cmd/worker/main.go`
  - `docs/v2/designs/2026-06-30-v2-tag-release-image.md`
- Checks:
  - `make build`
  - `docker build -t port-agent-worker:ci .` when Docker is available
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: Add GitHub Actions workflow for tag-triggered GHCR image publishing.
- Depends on:
  - T1
- Write Scope:
  - `.github/workflows/release-image.yml`
- Read Context:
  - `Dockerfile`
  - `docs/v2/designs/2026-06-30-v2-tag-release-image.md`
- Checks:
  - YAML parse check for `.github/workflows/release-image.yml`
- Parallel-safe: no

### Task T3

- [x] Complete
- Goal: Run repository verification for the release image workflow slice.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `Dockerfile`
  - `.dockerignore`
  - `.github/workflows/release-image.yml`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
  - `docker build -t port-agent-worker:ci .` when Docker is available
- Parallel-safe: no
