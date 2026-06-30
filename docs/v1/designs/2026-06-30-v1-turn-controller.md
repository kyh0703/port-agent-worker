---
feature: turn-controller
created_at: 2026-06-30T12:20:27+09:00
---

# Turn Controller

## Goal

음성 세션의 barge-in과 endpointing 판단을 application 계층의 순수 상태기계로 분리한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/ROADMAP.md`
  - `docs/v1/designs/2026-06-30-v1-pion-voice-worker.md`
- Existing system facts:
  - `internal/application/session`은 STT -> LLM -> TTS orchestration을 소유한다.
  - Pion, STT, LLM, TTS 구현은 adapter 계층에 격리되어 있다.
  - MVP에서는 Silero VAD와 LocalSmartTurnAnalyzerV3 실제 구현은 deferred이다.
- User brief:
  - Pion SFU를 유지한다.
  - DDD + Hexagonal Architecture를 유지한다.
  - BotStarted -> barge-in, BotStopped -> endpointing stop_secs 흐름이 필요하다.

## Plan Handoff

### Scope for Planning

- `internal/application/turn` 패키지를 추가한다.
- bot/user speech 이벤트를 받아 barge-in과 endpointing 결정을 반환한다.
- 상태기계는 외부 SDK 타입, Pion 타입, provider 타입을 참조하지 않는다.
- smart turn analyzer는 future adapter가 붙을 수 있는 application interface까지만 둔다.

### Success Criteria

- `BotStarted` 이후 사용자 발화 시작 시 barge-in 결정이 true가 된다.
- `BotStopped` 이후 사용자 발화가 없거나 멈춘 상태에서 `stop_secs`가 지나면 endpoint 결정이 true가 된다.
- 사용자 발화 중에는 endpoint가 발생하지 않는다.
- 기본 stop delay가 있고 잘못된 설정은 안전한 기본값으로 보정된다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- Silero VAD ONNX runtime integration.
- LocalSmartTurnAnalyzerV3 구현.
- session runner에 turn controller를 완전 통합.
- RTP track acquire/signaling.

### Open Questions

- 실제 VAD 이벤트 소스는 Pion media flow가 확정된 뒤 adapter에서 연결한다.

### Suggested Validation

- `go test ./internal/application/turn`
- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/application/turn/**`
- Shared files to avoid touching in parallel:
  - none expected.
- Likely sequential dependencies:
  - controller API -> tests.

