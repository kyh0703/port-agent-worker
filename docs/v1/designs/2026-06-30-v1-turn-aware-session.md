---
feature: turn-aware-session
created_at: 2026-06-30T12:31:28+09:00
---

# Turn Aware Session

## Goal

세션 orchestration에서 PCM stream을 STT와 VAD로 fan-out하고, VAD activity event를 turn decision으로 변환해 세션 경계에서 관찰할 수 있게 한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-voice-activity-flow.md`
  - `docs/v1/designs/2026-06-30-v1-turn-controller.md`
- Existing system facts:
  - `session.Orchestrator`는 ingress PCM을 STT에 넘기고 final transcript를 받으면 LLM/TTS/egress를 수행한다.
  - `ports.VoiceActivityDetector`는 PCM frame stream을 speech activity event stream으로 변환한다.
  - `turn.ActivityProcessor`는 speech activity event를 barge-in/endpoint decision으로 변환한다.
- User brief:
  - worker가 STT, LLM, TTS, 송출 역할을 한 프로세스에서 가져간다.
  - 내부는 DDD + Hexagonal Architecture로 쪼갠다.

## Plan Handoff

### Scope for Planning

- `internal/application/session`에 turn-aware runner 조립 함수를 추가한다.
- turn-aware orchestrator는 PCM input을 STT와 VAD로 fan-out한다.
- VAD 이벤트는 `turn.ActivityProcessor`로 처리하고, decision은 session-level handler로 전달한다.
- 기존 `NewRunnerFromRuntime`과 `Orchestrator` 기본 동작은 유지한다.

### Success Criteria

- turn runtime이 없는 기존 runner 테스트는 유지된다.
- turn-aware runner는 STT와 VAD가 같은 PCM frame stream을 받는다.
- speech started event가 barge-in decision으로 handler에 전달된다.
- speech stopped event 이후 endpoint decision이 handler에 전달된다.
- VAD/turn handler 오류는 context를 감싼 오류로 반환된다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- Silero VAD adapter 구현.
- TTS 송출 중단, egress interrupt 구현.
- Deepgram endpointing 옵션 변경.
- SFU signaling 또는 track acquire.

### Open Questions

- barge-in decision을 실제 TTS/egress interrupt로 연결하는 방식은 후속 feature에서 확정한다.

### Suggested Validation

- `go test ./internal/application/session`
- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/application/session/**`
- Shared files to avoid touching in parallel:
  - `internal/application/session/orchestrator.go`
  - `internal/application/session/assembler.go`
- Likely sequential dependencies:
  - turn runtime types -> orchestrator fan-out -> tests.

