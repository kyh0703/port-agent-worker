---
feature: pion-voice-worker
created_at: 2026-06-30T09:50:00+09:00
---

# Pion Voice Worker

## Goal

Pion SFU를 유지하면서 RTP audio track 기반 STT -> LLM -> TTS worker의 첫 실행 가능한 Go 골격을 만든다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/ROADMAP.md`
- Existing system facts:
  - 현재 repo는 새 프로젝트다.
  - SFU는 Pion 기반으로 유지한다.
  - worker는 SIP/WebRTC endpoint가 아니라 SFU track을 구독/발행한다.
- External constraints:
  - STT는 Deepgram nova-3를 목표 provider로 둔다.
  - LLM은 OpenRouter gemini-2.5-flash-lite를 목표 provider로 둔다.
  - TTS는 Cartesia/ElevenLabs를 목표 provider로 둔다.

## Problem Statement

Pion에서는 worker가 RTP packet을 직접 다루게 되므로 voice pipeline 앞뒤에 RTP <-> PCM media boundary가 필요하다. 제품 코어가 provider SDK와 media 세부 구현에 묶이면 barge-in, VAD, filler, RAG 확장이 어려워진다.

## Decision Drivers

- Pion SFU 유지.
- DDD + Hexagonal Architecture 적용.
- 초기 MVP는 실제 media codec보다 구조와 테스트 가능한 orchestration 우선.
- provider는 교체 가능해야 한다.
- GStreamer/LiveKit 전환 가능성을 adapter 경계로 남긴다.

## Options Considered

### Option A

- Summary: Pion + 순수 Go media adapter를 기본 방향으로 잡고, domain/application은 PCM 포트만 본다.
- Pros: 배포 단순, 기존 SFU 유지, adapter 교체 쉬움.
- Cons: 실제 Opus decode/encode와 jitter 처리는 v2에서 구현 부담.
- Risks: media 품질 이슈가 생기면 GStreamer fallback 필요.

### Option B

- Summary: GStreamer를 즉시 도입해 RTP/Opus/PCM 처리를 맡긴다.
- Pros: media pipeline 성숙도 높음.
- Cons: 설치형 native dependency, Docker/plugin/license 관리 증가.
- Risks: 초기 디버깅 포인트가 Go service 밖으로 늘어남.

## Recommended Option

- Choice: Option A.
- Why now: 첫 마일스톤은 voice worker의 경계와 orchestration을 검증하는 것이 중요하다.
- Rejected alternatives: GStreamer 즉시 도입은 MVP 대비 운영 복잡도가 크다.

## Scope Decision

- In:
  - Go module 초기화.
  - DDD + Hexagonal package layout.
  - audio/domain value object.
  - STT, LLM, TTS, audio ingress/egress 포트.
  - session orchestrator.
  - Pion RTP adapter skeleton.
  - fake provider 기반 테스트.
- Out:
  - 실제 Deepgram/OpenRouter/Cartesia/ElevenLabs API 호출.
  - 실제 Opus decode/encode/resample.
  - Silero VAD ONNX runtime.
  - RAG, filler, emotion, recording.
- Deferred:
  - GStreamer media adapter.
  - LocalSmartTurnAnalyzerV3.
  - barge-in output flush.

## Open Questions

- 실제 Pion SFU worker 접속 방식은 in-process hook인지, 별도 peer participant인지 v2에서 확정한다.
- TTS provider 우선순위는 Cartesia와 ElevenLabs 중 v2에서 확정한다.

## Plan Handoff

### Source of Truth Docs

- `docs/ARCHITECTURE.md`
- `docs/ROADMAP.md`
- `docs/v1/designs/2026-06-30-v1-pion-voice-worker.md`

### Scope for Planning

`pion-voice-worker` 단일 기능으로 Go worker skeleton과 테스트 가능한 session orchestration을 구현한다.

### Fixed Constraints

- Pion 유지.
- DDD + Hexagonal Architecture.
- domain/application 계층은 provider SDK와 Pion 타입을 import하지 않는다.
- 녹취는 제외한다.

### Success Criteria

- `go test ./...` 통과.
- `go build ./cmd/worker` 통과.
- fake ingress/STT/LLM/TTS/egress로 사용자 transcript가 assistant audio write까지 흐르는 테스트가 있다.
- Pion adapter는 application port를 구현하는 형태로 컴파일된다.

### Non-Goals

- Production media quality.
- 실제 외부 provider API 호출.
- SFU signaling integration.
- 배포 자동화.

### Open Questions

- 실제 SFU와의 접속 방식.
- production codec/resampling 구현 선택.

### Suggested Validation

- `go test ./...`
- `go build ./cmd/worker`

### Parallelization Hints

- Candidate write boundaries:
  - domain/application core.
  - adapter skeleton.
  - command wiring/config.
- Shared files to avoid touching in parallel:
  - `go.mod`
  - `internal/application/session/orchestrator.go`
- Likely sequential dependencies:
  - domain values -> ports -> orchestrator -> adapters -> cmd wiring.

