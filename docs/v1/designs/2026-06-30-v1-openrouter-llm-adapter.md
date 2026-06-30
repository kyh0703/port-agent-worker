---
feature: openrouter-llm-adapter
created_at: 2026-06-30T12:05:00+09:00
---

# OpenRouter LLM Adapter

## Goal

`LanguageModel` 포트 뒤에 OpenRouter chat completions adapter를 추가한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-deepgram-stt-adapter.md`
  - `docs/v1/designs/2026-06-30-v1-cartesia-tts-adapter.md`
- Existing system facts:
  - application port는 `LanguageModel.Generate(ctx, voice.UserUtterance)`이다.
  - 현재 포트는 non-streaming response를 반환한다.
  - provider SDK/HTTP 모델은 adapter 내부에 둔다.
- User brief:
  - LLM은 OpenRouter `google/gemini-2.5-flash-lite`를 사용한다.
  - DDD + Hexagonal Architecture를 유지한다.

## Plan Handoff

### Scope for Planning

- `internal/adapters/providers/openrouter` 패키지를 추가한다.
- OpenRouter `/api/v1/chat/completions`를 호출한다.
- user utterance와 optional system prompt를 chat messages로 변환한다.
- 응답 `choices[0].message.content`를 `voice.AssistantResponse`로 변환한다.

### Success Criteria

- OpenRouter adapter가 `ports.LanguageModel`을 구현한다.
- request header에 `Authorization: Bearer`가 포함된다.
- request body에 `model=google/gemini-2.5-flash-lite`와 messages가 포함된다.
- response parsing 테스트와 HTTP error 테스트가 있다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- streaming token output.
- tool calling.
- reasoning parameter tuning.
- RAG context injection.
- 실제 OpenRouter API 통합 테스트.

### Open Questions

- system prompt는 운영 설정에서 주입한다. 기본값은 비워둔다.

### Suggested Validation

- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/providers/openrouter/**`
- Shared files to avoid touching in parallel:
  - none expected.
- Likely sequential dependencies:
  - request/response model -> HTTP client -> tests.

