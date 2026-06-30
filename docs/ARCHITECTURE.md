# Architecture

## Purpose

Pion SFU 위에서 동작하는 voice agent worker를 DDD + Hexagonal Architecture로 구성한다. 코어는 음성 세션과 턴 처리 규칙을 소유하고, WebRTC/RTP, STT, LLM, TTS는 어댑터로 격리한다.

## Shared Boundaries

- Domain:
  - `internal/domain/voice`: 세션, PCM frame, transcript, assistant response, turn state.
- Application:
  - `internal/application/ports`: ingress/egress, STT, LLM, TTS, VAD 포트.
  - `internal/application/session`: STT -> LLM -> TTS orchestration.
- Adapters:
  - `internal/adapters/pionrtp`: Pion RTP ingress/egress.
  - `internal/adapters/providers`: provider 구현.
- Entrypoint:
  - `cmd/worker`: config 로딩, signal handling, dependency wiring.

## Shared Constraints

- Security:
  - Provider API key는 환경변수로만 읽는다.
  - domain/application 계층에 secret을 전달하지 않는다.
- Reliability:
  - 모든 long-running loop는 `context.Context` 취소를 따른다.
  - provider 오류는 wrapping해서 세션 경계에서 처리한다.
- Performance:
  - audio frame은 작은 immutable value로 다룬다.
  - MVP는 backpressure를 channel buffer와 context cancellation으로 제어한다.
- Operational limits:
  - Pion media adapter는 MVP에서 skeleton/mock 중심이다.
  - 실제 codec/resampling 최적화는 v2에서 확정한다.

