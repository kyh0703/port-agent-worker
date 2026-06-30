# Pion Voice Worker

## Workspace

- Branch: feature/v1-pion-voice-worker
- Base: main
- Isolation: `.worktrees/feature-v1-pion-voice-worker`

## Goal

Pion RTP 기반 voice agent worker의 DDD + Hexagonal skeleton을 만들고, fake provider 테스트로 STT -> LLM -> TTS -> egress 흐름을 검증한다.

## Source

- Design: `docs/v1/designs/2026-06-30-v1-pion-voice-worker.md`
- Architecture: `docs/ARCHITECTURE.md`

## Task Graph

### T1 Domain and Ports

- [ ] Complete
- Goal: audio/session domain value와 application port를 정의한다.
- Depends on: none
- Write Scope: `internal/domain/**`, `internal/application/ports/**`
- Read Context: `docs/ARCHITECTURE.md`
- Checks: `go test ./...`
- Parallel-safe: no

### T2 Session Orchestrator

- [ ] Complete
- Goal: ingress -> STT -> LLM -> TTS -> egress orchestration을 구현하고 fake 테스트를 작성한다.
- Depends on: T1
- Write Scope: `internal/application/session/**`
- Read Context: `internal/application/ports/**`, `internal/domain/**`
- Checks: `go test ./...`
- Parallel-safe: no

### T3 Adapter Skeleton and Command Wiring

- [ ] Complete
- Goal: Pion RTP adapter skeleton, noop provider, config, worker entrypoint를 구성한다.
- Depends on: T1, T2
- Write Scope: `internal/adapters/**`, `internal/config/**`, `cmd/worker/**`
- Read Context: `internal/application/ports/**`
- Checks: `go test ./...`, `go build ./cmd/worker`
- Parallel-safe: no

### T4 Project Tooling

- [ ] Complete
- Goal: Makefile, README, module metadata를 정리한다.
- Depends on: T3
- Write Scope: `Makefile`, `README.md`, `go.mod`
- Read Context: full tree
- Checks: `make test`, `make build`
- Parallel-safe: no

