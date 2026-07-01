---
feature: rename-port-voice-pipeline
created_at: 2026-07-01T21:17:28+09:00
---

# Rename Port Voice Pipeline

## Goal

Rename the module and runtime identity from `port-agent-worker` to `port-voice-pipeline` so the repository name matches its voice pipeline responsibility.

## Context / Inputs

- Source docs:
  - `docs/STATE.md`
  - `docs/ROADMAP.md`
  - `docs/ARCHITECTURE.md`
- Existing system facts:
  - The service currently owns RTP/PCM ingress, STT, LLM, TTS, RTP/PCM egress, and turn/VAD orchestration.
  - `cmd/worker` is still the current Go entrypoint path.
  - Build, Docker, release, and logs still use `port-agent-worker`.
- User brief:
  - Rename option 1 to `port-voice-pipeline` and update docs accordingly.

## Plan Handoff

### Scope for Planning

- Change Go module path and internal imports from `port-agent-worker` to `port-voice-pipeline`.
- Change binary, GHCR image identity, Docker image tag examples, Docker entrypoint, deployment commands, logs, app title, and README heading to `port-voice-pipeline`.
- Update architecture/state/roadmap/current v2 docs to describe the service as a voice pipeline instead of an agent worker where the wording refers to product identity.
- Keep `cmd/worker` path unchanged unless a plan proves a directory rename is required.

### Success Criteria

- `go test ./...` passes.
- `make build` produces `bin/port-voice-pipeline`.
- Repository text no longer uses `port-agent-worker` as the active module/runtime identity.
- Historical docs may keep `cmd/worker` path references where they describe the existing entrypoint.

### Non-Goals

- Do not implement gateway dispatch.
- Do not change runtime behavior.
- Do not rename `cmd/worker` unless needed for compile/build correctness.
- Do not change provider, media, turn, or domain contracts.

### Open Questions

- None.

### Suggested Validation

- `go test ./...`
- `make build`
- `rg -n "port-agent-worker|voice agent worker|agent worker" README.md docs go.mod Makefile Dockerfile rails.toml .github cmd internal`

### Parallelization Hints

- Candidate write boundaries:
  - `go.mod`
  - `cmd/**`
  - `internal/**`
  - `.github/workflows/**`
  - `README.md`
  - `Makefile`
  - `Dockerfile`
  - `rails.toml`
  - `docs/**`
- Shared files to avoid touching in parallel:
  - `go.mod`
  - `docs/STATE.md`
  - `docs/ARCHITECTURE.md`
- Likely sequential dependencies:
  - Module path/import rename first, then runtime/documentation naming updates, then verification.
