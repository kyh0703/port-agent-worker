---
feature: tag-release-image
created_at: 2026-06-30T16:10:55+09:00
---

# Tag Release Image

## Goal

GitHub tag push 시 worker Docker image를 빌드하고 GitHub Container Registry(GHCR)에 push한다.

## Context / Inputs

- Source docs:
  - `docs/STATE.md`
  - `docs/ROADMAP.md`
  - `docs/ARCHITECTURE.md`
- Existing system facts:
  - Go voice pipeline entrypoint is `./cmd/worker`.
  - Current build command is `go build -o bin/port-voice-pipeline ./cmd/worker`.
  - No Dockerfile or GitHub Actions workflow exists yet.
  - Deployment config exists in `rails.toml`.
- User brief:
  - "github에 tag딸떄 이미지 따서 github에 이미지 올릴려고 하거든 github action구현"

## Plan Handoff

### Scope for Planning

- Add a Dockerfile for the Go voice pipeline image.
- Add a GitHub Actions workflow that runs on semantic version tag pushes and pushes image tags to GHCR.
- Use repository-local build/test checks before pushing the image.

### Success Criteria

- Pushing a tag like `v1.2.3` triggers the workflow.
- Workflow authenticates to `ghcr.io` using `GITHUB_TOKEN`.
- Workflow builds the voice pipeline image and pushes it to `ghcr.io/<owner>/port-voice-pipeline`.
- Workflow publishes tag-derived image tags and labels.
- Local verification confirms workflow YAML and Dockerfile are syntactically usable.

### Non-Goals

- Do not publish images on every branch push.
- Do not add separate cloud provider deployment.
- Do not add runtime secrets to the repository.
- Do not change application runtime behavior.

### Open Questions

- None.

### Suggested Validation

- `go test ./...`
- `make build`
- `docker build -t port-voice-pipeline:ci .` when Docker is available
- YAML parse check for `.github/workflows/release-image.yml`

### Parallelization Hints

- Candidate write boundaries:
  - `Dockerfile`
  - `.dockerignore`
  - `.github/workflows/release-image.yml`
- Shared files to avoid touching in parallel:
  - none
- Likely sequential dependencies:
  - Dockerfile first, then workflow.
