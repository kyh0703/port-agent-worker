---
feature: silero-vad-adapter
created_at: 2026-06-30T12:51:09+09:00
---

# Silero VAD Adapter

## Goal

Silero VAD ONNX 런타임을 나중에 붙일 수 있도록 `ports.VoiceActivityDetector` adapter skeleton과 speech/silence 상태 전이를 구현한다.

## Context / Inputs

- Source docs:
  - `docs/ARCHITECTURE.md`
  - `docs/v1/designs/2026-06-30-v1-voice-activity-flow.md`
  - `docs/v1/designs/2026-06-30-v1-turn-runtime-wiring.md`
- Existing system facts:
  - `ports.VoiceActivityDetector`는 PCM frame stream을 `voice.SpeechActivityEvent` stream으로 변환한다.
  - `internal/adapters/vad/noop`은 frame drain만 수행한다.
  - 실제 Silero ONNX runtime integration은 roadmap에서 deferred이다.
- User brief:
  - Pion SFU 위에서 STT -> LLM -> TTS worker 구조를 유지한다.
  - turn-taking은 smart turn off 시 VAD 묵음 기반으로 동작해야 한다.

## Plan Handoff

### Scope for Planning

- `internal/adapters/vad/silero` 패키지를 추가한다.
- adapter는 `Engine` interface를 받아 frame별 speech probability를 얻는다.
- `speech threshold`, `min speech frames`, `min silence frames` 설정으로 started/stopped event를 만든다.
- ONNX runtime, model loading, resampling은 이번 범위에 넣지 않는다.

### Success Criteria

- `silero.Detector`는 `ports.VoiceActivityDetector`를 구현한다.
- speech probability가 threshold 이상으로 충분히 지속되면 `SpeechStarted` event를 낸다.
- silence probability가 충분히 지속되면 `SpeechStopped` event를 낸다.
- 입력 frame stream이 닫히면 event stream도 닫힌다.
- context 취소 시 event stream이 닫힌다.
- engine 누락/잘못된 config는 명시적 에러를 반환한다.
- `go test ./...`와 `make build`가 통과한다.

### Non-Goals

- ONNX runtime dependency 추가.
- Silero model 파일 다운로드/로딩.
- resampling implementation.
- turn runtime factory의 기본 VAD를 Silero로 교체.

### Suggested Validation

- `go test ./internal/adapters/vad/silero`
- `go test ./...`
- `make build`

### Parallelization Hints

- Candidate write boundaries:
  - `internal/adapters/vad/silero/**`
- Shared files to avoid touching in parallel:
  - none expected.
- Likely sequential dependencies:
  - detector API -> state transition tests -> implementation.

