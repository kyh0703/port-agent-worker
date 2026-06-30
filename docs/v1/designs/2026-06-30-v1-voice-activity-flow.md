---
feature: voice-activity-flow
created_at: 2026-06-30T12:26:17+09:00
---

# Voice Activity Flow

## Goal

VAD/Silero adapter가 나중에 붙을 수 있도록 음성 활동 이벤트와 turn controller 연결 계약을 DDD + Hexagonal 경계 안에 추가한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-turn-controller.md`
- Existing system facts:
  - `internal/application/turn.Controller`는 bot/user speech 상태와 endpointing 판단을 소유한다.
  - `internal/application/ports`에는 audio/provider 포트만 있고 VAD 포트는 아직 없다.
  - Silero VAD 실제 구현은 MVP deferred이다.
- User brief:
  - Pion SFU RTP를 유지한다.
  - STT -> LLM -> TTS worker가 turn-taking을 가져야 한다.
  - smart turn off 상태에서는 VAD 묵음만으로 endpointing해야 한다.

## Plan Handoff

### Scope for Planning

- `internal/domain/voice`에 speech activity event value를 추가한다.
- `internal/application/ports`에 VAD adapter용 포트를 추가한다.
- `internal/application/turn`에 activity event를 controller method로 변환하는 processor를 추가한다.
- 외부 SDK 타입과 Pion 타입은 domain/application 경계를 넘지 않는다.

### Success Criteria

- speech started/stopped 이벤트가 domain value로 표현된다.
- VAD adapter가 PCM frame을 받아 speech activity event stream을 반환하는 port가 있다.
- turn processor가 speech started를 barge-in 결정으로 변환한다.
- turn processor가 speech stopped 뒤 stop delay 경과를 endpoint 결정으로 변환한다.
- 잘못된 activity event는 명시적으로 에러 처리된다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- Silero VAD 구현.
- session runner에 processor 완전 통합.
- TTS 송출 중단 구현.
- RTP track acquire/signaling.

### Open Questions

- 실제 VAD adapter는 Silero ONNX 또는 플랫폼 VAD 중 하나로 후속 feature에서 선택한다.

### Suggested Validation

- `go test ./internal/domain/voice`
- `go test ./internal/application/ports`
- `go test ./internal/application/turn`
- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/domain/voice/**`
  - `internal/application/ports/**`
  - `internal/application/turn/**`
- Shared files to avoid touching in parallel:
  - `internal/application/turn/controller.go`
- Likely sequential dependencies:
  - domain event -> VAD port -> turn processor.

