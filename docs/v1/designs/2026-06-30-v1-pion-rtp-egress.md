---
feature: pion-rtp-egress
created_at: 2026-06-30T13:05:00+09:00
---

# Pion RTP Egress

## Goal

PCM frame stream을 Opus RTP packet으로 packetize해서 Pion output track에 쓰는 egress adapter 경계를 구현한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-pion-rtp-ingress.md`
- Existing system facts:
  - `AudioEgress` 포트는 `WritePCM(ctx, voice.PCMFrame)`이다.
  - `internal/adapters/pionrtp.Egress`는 현재 `ErrMediaPipelineNotReady`만 반환한다.
  - pure Go `pion/opus`는 decoder만 제공한다.
  - 현재 환경에는 `libopus` pkg-config가 없다.
- User brief:
  - Pion SFU를 유지한다.
  - DDD + Hexagonal Architecture를 유지한다.
  - worker가 TTS audio를 SFU RTP track으로 송출해야 한다.

## Plan Handoff

### Scope for Planning

- `pionrtp.Egress`에 encoder 인터페이스와 RTP writer 인터페이스를 추가한다.
- PCM frame을 encoder에 전달하고 encoded Opus payload를 RTP packet으로 작성한다.
- sequence number, RTP timestamp, payload type, SSRC를 adapter config로 관리한다.
- production encoder 구현은 제외하고 fake encoder 기반 테스트로 packetization 경계를 검증한다.

### Success Criteria

- `Egress.WritePCM`이 encoder가 있는 경우 더 이상 `ErrMediaPipelineNotReady`를 반환하지 않는다.
- domain/application 계층에는 Pion/codec 타입이 노출되지 않는다.
- fake encoder/writer 기반 테스트가 있다.
- RTP timestamp가 PCM frame duration 기준으로 증가한다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- production Opus encoder 선택.
- libopus 설치/Docker 작업.
- GStreamer adapter.
- jitter/timing scheduler.
- 실제 SFU 통합 테스트.

### Open Questions

- production encoder는 v2에서 `libopus` CGO, GStreamer, 별도 media service 중 하나로 확정한다.

### Suggested Validation

- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/pionrtp/**`
- Shared files to avoid touching in parallel:
  - none expected.
- Likely sequential dependencies:
  - encoder/writer boundary -> RTP packetizer -> tests.

