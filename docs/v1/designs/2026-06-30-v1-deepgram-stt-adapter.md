---
feature: deepgram-stt-adapter
created_at: 2026-06-30T11:20:00+09:00
---

# Deepgram STT Adapter

## Goal

`SpeechToText` 포트 뒤에 Deepgram Nova-3 streaming STT adapter를 추가한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/plans/pion-voice-worker.md`
- Existing system facts:
  - application port는 `SpeechToText.Transcribe(ctx, <-chan voice.PCMFrame)`이다.
  - domain/application 계층은 provider SDK나 transport 타입을 import하면 안 된다.
  - 현재 Pion media adapter는 아직 실제 PCM 생성 전 skeleton이다.
- User brief:
  - STT/TTS/LLM 연결을 순차적으로 진행한다.
  - DDD + Hexagonal Architecture를 유지한다.
  - STT는 Deepgram Nova-3를 사용한다.

## Plan Handoff

### Scope for Planning

- `internal/adapters/providers/deepgram` 패키지를 추가한다.
- PCMFrame stream을 Deepgram WebSocket streaming request로 전송한다.
- Deepgram transcript JSON을 `voice.Transcript`로 변환한다.
- provider 설정은 adapter config 구조체로 받는다.
- application/domain 경계는 유지한다.

### Success Criteria

- Deepgram adapter가 `ports.SpeechToText`를 구현한다.
- `linear16`, `sample_rate`, `channels`, `model=nova-3`, `language` query가 구성된다.
- final/interim transcript JSON parsing 테스트가 있다.
- audio send loop가 PCM frame data를 WebSocket writer에 전달한다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- 실제 Deepgram API 통합 테스트.
- Pion RTP -> PCM decode 구현.
- endpointing, VAD, smart turn 정책 구현.
- retry/backoff/metrics 구현.

### Open Questions

- Deepgram language 기본값은 우선 `ko`로 둔다.
- endpointing 기본값은 provider default를 사용하고, turn controller 구현 시 재검토한다.

### Suggested Validation

- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/providers/deepgram/**`
  - `internal/config/**`
- Shared files to avoid touching in parallel:
  - `go.mod`
  - `go.sum`
- Likely sequential dependencies:
  - adapter config -> WebSocket URL builder/parser tests -> Transcribe implementation.

