---
feature: pion-media-runtime
created_at: 2026-06-30T15:00:00+09:00
---

# Pion Media Runtime

## Goal

이미 획득한 Pion input/output track을 session `AudioRuntime`으로 변환하는 media factory를 추가한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-pion-rtp-ingress.md`
  - `docs/v1/designs/2026-06-30-v1-pion-rtp-egress.md`
  - `docs/v1/designs/2026-06-30-v1-worker-session-wiring.md`
- Existing system facts:
  - `pionrtp.NewIngress(*webrtc.TrackRemote)`가 RTP -> PCM ingress를 만든다.
  - `pionrtp.NewEgressWithEncoder(*webrtc.TrackLocalStaticRTP, encoder)`가 PCM -> RTP egress boundary를 만든다.
  - `session.AudioRuntime`은 `ports.AudioIngress`/`ports.AudioEgress`를 받는다.
  - 실제 SFU signaling/track acquisition은 아직 없다.
- User brief:
  - Pion SFU를 유지한다.
  - DDD + Hexagonal Architecture를 유지한다.

## Plan Handoff

### Scope for Planning

- `pionrtp.FrameEncoder` 인터페이스를 public으로 정리한다.
- `internal/adapters/media/pion` 패키지를 추가한다.
- input/output Pion track과 optional encoder를 `session.AudioRuntime`으로 변환한다.
- track/encoder 누락 시 명시적 에러를 반환한다.

### Success Criteria

- `media/pion.NewRuntime`이 `session.AudioRuntime`을 반환한다.
- Pion/provider 타입은 adapter 계층 밖으로 확산되지 않는다.
- egress production encoder가 없으면 명시적 config error로 실패한다.
- fake encoder 기반 factory 테스트가 있다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- SFU signaling.
- PeerConnection participant 구현.
- production Opus encoder.
- 실제 Pion integration test.

### Open Questions

- SFU track acquisition 방식은 다음 feature에서 확정한다.

### Suggested Validation

- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/pionrtp/**`
  - `internal/adapters/media/pion/**`
- Shared files to avoid touching in parallel:
  - none expected.
- Likely sequential dependencies:
  - exported encoder interface -> media runtime factory -> tests.

