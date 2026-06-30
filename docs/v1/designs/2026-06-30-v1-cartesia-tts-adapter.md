---
feature: cartesia-tts-adapter
created_at: 2026-06-30T11:45:00+09:00
---

# Cartesia TTS Adapter

## Goal

`TextToSpeech` 포트 뒤에 Cartesia Sonic TTS adapter를 추가한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-deepgram-stt-adapter.md`
- Existing system facts:
  - application port는 `TextToSpeech.Synthesize(ctx, voice.AssistantResponse)`이다.
  - 현재 LLM token streaming 포트가 없으므로 complete transcript 기반 TTS가 자연스럽다.
  - provider SDK 타입은 adapter 밖으로 노출하지 않는다.
- User brief:
  - TTS는 Cartesia를 사용한다.
  - DDD + Hexagonal Architecture를 유지한다.

## Plan Handoff

### Scope for Planning

- `internal/adapters/providers/cartesia` 패키지를 추가한다.
- Cartesia `/tts/bytes` API를 호출해 raw audio bytes를 받는다.
- response body를 `voice.PCMFrame` stream으로 chunking한다.
- model, voice id, language, output format, API version을 adapter config로 둔다.

### Success Criteria

- Cartesia adapter가 `ports.TextToSpeech`를 구현한다.
- request header에 API key와 `Cartesia-Version`이 포함된다.
- request body에 `model_id`, `transcript`, `voice`, `output_format`, `language`가 포함된다.
- response audio bytes가 `voice.PCMFrame`으로 방출되는 테스트가 있다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- Cartesia WebSocket input streaming.
- word timestamp 처리.
- emotion/generation_config 세부 제어.
- 실제 Cartesia API 통합 테스트.

### Open Questions

- voice id는 운영 설정에서 받아야 한다. 기본값은 두지 않고 필수값으로 처리한다.

### Suggested Validation

- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/providers/cartesia/**`
- Shared files to avoid touching in parallel:
  - `go.mod`
  - `go.sum`
- Likely sequential dependencies:
  - request model/config -> HTTP client -> chunking tests.

