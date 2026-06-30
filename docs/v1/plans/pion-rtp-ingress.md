# Pion RTP Ingress

## Goal

Pion audio track의 RTP packet을 PCM frame stream으로 변환하는 `AudioIngress` adapter를 구현한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-pion-rtp-ingress.md`

## Workspace

- Branch: feature/v1-pion-rtp-ingress
- Base: main
- Isolation: `.worktrees/feature-v1-pion-rtp-ingress`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [ ] Complete
- Goal: Pion packet source와 Opus decoder adapter boundary를 만든다.
- Depends on:
  - none
- Write Scope:
  - `internal/adapters/pionrtp/**`
- Read Context:
  - `internal/adapters/pionrtp/ingress.go`
  - `internal/application/ports/audio.go`
  - `internal/domain/voice/audio.go`
- Checks:
  - `go test ./...`
- Parallel-safe: no

### Task T2

- [ ] Complete
- Goal: `Ingress.PCMFrames` loop를 구현하고 fake 기반 테스트를 추가한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/pionrtp/**`
- Read Context:
  - `internal/domain/voice/audio.go`
- Checks:
  - `go test ./internal/adapters/pionrtp`
- Parallel-safe: no

### Task T3

- [ ] Complete
- Goal: dependency 정리와 전체 검증을 완료한다.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `go.mod`
  - `go.sum`
  - `internal/adapters/pionrtp/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no

