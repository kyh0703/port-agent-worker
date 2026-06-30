---
feature: pipeline-stability-hex-wiring
created_at: 2026-06-30T13:28:52+09:00
---

# Pipeline Stability and Hexagonal Wiring Cleanup

## Goal

Turn-aware audio fan-out이 한 consumer 종료로 멈추지 않게 하고, provider wiring에서 adapter package가 global config를 직접 알지 않게 한다.

## Context / Inputs

- Source docs:
  - `docs/STATE.md`
  - `docs/ROADMAP.md`
  - `docs/ARCHITECTURE.md`
- Existing system facts:
  - `internal/application/session` owns STT -> LLM -> TTS orchestration and turn-aware audio fan-out.
  - `internal/application/ports` defines provider/audio/VAD boundaries.
  - `internal/adapters/providers` builds Deepgram/OpenRouter/Cartesia clients.
  - `cmd/worker` is the current composition root.
- Bug brief:
  - `TurnAwareOrchestrator` can block when one fan-out consumer stops reading before the other.
  - `internal/adapters/providers` imports `internal/config`, which weakens Hexagonal composition boundaries.
- Reproduction:
  - Use a VAD that returns/closes while STT still waits for audio channel close; the unbuffered `vadAudio` send can block `fanOutPCM`, preventing STT input closure.
  - Inspect `internal/adapters/providers/factory.go`; it imports `internal/config` and translates environment-level config inside the adapter package.
- Expected vs Actual:
  - Expected: a stopped VAD consumer should not block STT audio delivery/closure.
  - Actual: fan-out sends sequentially to all outputs and can block on a consumer no longer reading.
  - Expected: composition root translates `config.Config` to adapter-specific configs.
  - Actual: provider adapter package knows `config.Config`.
- Suspected Cause:
  - Fan-out does not detach or drop a closed/stalled output.
  - Provider runtime factory combines bootstrap config mapping with adapter construction.
- Regression Risk:
  - Audio fan-out changes can alter cancellation/backpressure behavior.
  - Provider factory API changes can break worker wiring and tests.

## Plan Handoff

### Scope for Planning

- Add a focused test for turn-aware fan-out where VAD exits before STT finishes consuming input.
- Modify fan-out behavior so one blocked consumer cannot stall the remaining pipeline.
- Move `config.Config` -> provider-specific config mapping out of `internal/adapters/providers`.
- Preserve provider adapter implementations and public application ports.

### Success Criteria

- `TurnAwareOrchestrator` completes when VAD stops reading before STT has consumed all audio.
- `internal/adapters/providers` no longer imports `internal/config`.
- `cmd/worker` or a composition-level helper owns environment config translation.
- `go test ./...`, `go vet ./...`, and `make build` pass.

### Non-Goals

- Do not implement actual Pion track acquisition.
- Do not add new providers.
- Do not change domain/application port method signatures unless required for the fix.
- Do not implement Silero ONNX runtime.

### Open Questions

- None.

### Suggested Validation

- `go test ./internal/application/session`
- `go test ./internal/adapters/providers ./cmd/worker`
- `go test ./...`
- `go vet ./...`
- `make build`
- `go list -f '{{.ImportPath}} -> {{join .Imports " "}}' ./internal/adapters/providers` shows no `internal/config` import.

### Parallelization Hints

- Candidate write boundaries:
  - `internal/application/session/**`
  - `internal/adapters/providers/**`
  - `cmd/worker/**`
- Shared files to avoid touching in parallel:
  - `cmd/worker/main.go`
  - `internal/adapters/providers/factory.go`
- Likely sequential dependencies:
  - Fix fan-out first, then provider wiring cleanup, then full verification.
