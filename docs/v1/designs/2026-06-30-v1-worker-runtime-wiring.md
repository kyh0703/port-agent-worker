---
feature: worker-runtime-wiring
created_at: 2026-06-30T13:30:00+09:00
---

# Worker Runtime Wiring

## Goal

`cmd/worker`에서 환경 설정을 읽어 provider adapter와 session orchestrator를 조립할 수 있게 한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-deepgram-stt-adapter.md`
  - `docs/v1/designs/2026-06-30-v1-cartesia-tts-adapter.md`
  - `docs/v1/designs/2026-06-30-v1-openrouter-llm-adapter.md`
- Existing system facts:
  - `cmd/worker`는 현재 config를 로그로 출력하고 대기만 한다.
  - Deepgram, OpenRouter, Cartesia adapter는 이미 포트를 구현한다.
  - Pion track 획득/signaling은 아직 없다.
  - Pion egress production Opus encoder는 아직 없다.
- User brief:
  - STT, LLM, TTS, Pion media를 순차적으로 붙인다.
  - DDD + Hexagonal Architecture를 유지한다.

## Plan Handoff

### Scope for Planning

- config에 provider별 runtime 설정을 추가한다.
- provider factory를 adapter 계층에 추가해 config -> STT/LLM/TTS 포트 구현체로 변환한다.
- worker main이 provider factory를 호출해 wiring 가능 상태를 검증한다.
- 실제 Pion track 연결 전까지 ingress/egress는 명시적 pending 상태로 둔다.

### Success Criteria

- config가 Deepgram/OpenRouter/Cartesia 필수 env를 로드한다.
- provider factory가 `ports.SpeechToText`, `ports.LanguageModel`, `ports.TextToSpeech` 구현체를 만든다.
- main은 provider wiring 성공/실패를 명확히 로그/exit code로 처리한다.
- domain/application 계층에 provider config 타입이 노출되지 않는다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- Pion track acquisition/signaling.
- 실제 session 실행.
- production Opus encoder.
- provider API integration test.

### Open Questions

- Pion SFU와 worker 연결 방식은 별도 feature에서 확정한다.

### Suggested Validation

- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/config/**`
  - `internal/adapters/providers/**`
  - `cmd/worker/**`
- Shared files to avoid touching in parallel:
  - `cmd/worker/main.go`
- Likely sequential dependencies:
  - config -> provider factory -> main wiring.

