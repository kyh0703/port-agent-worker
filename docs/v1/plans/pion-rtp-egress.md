# Pion RTP Egress

## Goal

PCM frame을 encoded Opus RTP packet으로 packetize해 Pion output track에 쓰는 `AudioEgress` adapter 경계를 구현한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-pion-rtp-egress.md`

## Workspace

- Branch: feature/v1-pion-rtp-egress
- Base: main
- Isolation: `.worktrees/feature-v1-pion-rtp-egress`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [x] Complete
- Goal: Egress encoder/writer boundary와 config를 만든다.
- Depends on:
  - none
- Write Scope:
  - `internal/adapters/pionrtp/**`
- Read Context:
  - `internal/adapters/pionrtp/egress.go`
  - `internal/application/ports/audio.go`
  - `internal/domain/voice/audio.go`
- Checks:
  - `go test ./...`
- Parallel-safe: no

### Task T2

- [x] Complete
- Goal: RTP packetization과 fake 기반 테스트를 구현한다.
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

- [x] Complete
- Goal: 전체 검증을 완료한다.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `internal/adapters/pionrtp/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no
