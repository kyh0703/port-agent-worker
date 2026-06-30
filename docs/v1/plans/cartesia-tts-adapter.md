# Cartesia TTS Adapter

## Goal

Cartesia `/tts/bytes` 기반 TTS adapter를 추가해서 assistant response text를 PCM frame stream으로 변환한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-cartesia-tts-adapter.md`
- Cartesia Text-to-Speech Bytes docs.

## Workspace

- Branch: feature/v1-cartesia-tts-adapter
- Base: main
- Isolation: `.worktrees/feature-v1-cartesia-tts-adapter`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [ ] Complete
- Goal: Cartesia adapter config와 request body 모델을 구현한다.
- Depends on:
  - none
- Write Scope:
  - `internal/adapters/providers/cartesia/**`
- Read Context:
  - `internal/application/ports/providers.go`
  - `internal/domain/voice/**`
  - `docs/v1/designs/2026-06-30-v1-cartesia-tts-adapter.md`
- Checks:
  - `go test ./...`
- Parallel-safe: no

### Task T2

- [ ] Complete
- Goal: HTTP `/tts/bytes` 호출과 audio chunk -> PCMFrame 변환을 구현한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/providers/cartesia/**`
- Read Context:
  - `internal/domain/voice/audio.go`
- Checks:
  - `go test ./internal/adapters/providers/cartesia`
- Parallel-safe: no

### Task T3

- [ ] Complete
- Goal: adapter 테스트와 전체 빌드 검증을 완료한다.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `internal/adapters/providers/cartesia/**`
- Read Context:
  - `internal/application/ports/providers.go`
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no

