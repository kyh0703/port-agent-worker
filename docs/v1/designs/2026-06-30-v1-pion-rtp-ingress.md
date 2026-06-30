---
feature: pion-rtp-ingress
created_at: 2026-06-30T12:30:00+09:00
---

# Pion RTP Ingress

## Goal

Pion audio track에서 들어오는 Opus RTP packet을 PCM frame stream으로 변환하는 ingress adapter를 구현한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-pion-voice-worker.md`
- Existing system facts:
  - `AudioIngress` 포트는 `PCMFrames(ctx) <-chan voice.PCMFrame`이다.
  - `internal/adapters/pionrtp.Ingress`는 현재 `ErrMediaPipelineNotReady`만 반환한다.
  - provider adapter들은 16kHz mono PCM을 기대한다.
- User brief:
  - Pion SFU를 유지한다.
  - worker는 SFU의 RTP media track을 받아 STT -> LLM -> TTS 흐름으로 처리한다.
  - DDD + Hexagonal Architecture를 유지한다.

## Plan Handoff

### Scope for Planning

- `pionrtp.Ingress`를 실제 packet read loop로 바꾼다.
- Pion `TrackRemote.ReadRTP()` 결과를 adapter 내부 packet source로 추상화한다.
- Opus RTP payload를 PCM으로 decode하는 decoder boundary를 둔다.
- decoded PCM을 16kHz mono `voice.PCMFrame`으로 방출한다.
- 테스트는 fake packet source/decoder로 packet -> frame flow, nil track, cancellation을 검증한다.

### Success Criteria

- `pionrtp.Ingress`가 더 이상 정상 track에서 `ErrMediaPipelineNotReady`를 반환하지 않는다.
- domain/application 계층에 Pion 타입이나 codec 타입이 노출되지 않는다.
- packet source와 decoder가 adapter 내부 인터페이스로 분리된다.
- fake packet source/decoder 기반 테스트가 있다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- RTP jitter buffer 고도화.
- packet loss concealment.
- PCM -> RTP egress.
- resampler 품질 튜닝.
- 실제 SFU 통합 테스트.

### Open Questions

- production resampling 품질 요구가 생기면 v2에서 GStreamer 또는 전용 resampler를 비교한다.

### Suggested Validation

- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/pionrtp/**`
- Shared files to avoid touching in parallel:
  - `go.mod`
  - `go.sum`
- Likely sequential dependencies:
  - packet source abstraction -> decoder abstraction -> Ingress loop -> tests.

