# Architecture

## Purpose

Pion SFU 위에서 동작하는 voice pipeline runtime을 DDD + Hexagonal Architecture로 구성한다. 코어는 음성 세션과 턴 처리 규칙을 소유하고, WebRTC/RTP, STT, LLM, TTS는 어댑터로 격리한다.

`port-voice-pipeline`은 public API가 아니다. `../port-gateway`가
`../port-api` dispatch를 받은 뒤 agent participant attach를 지시하는
내부 voice pipeline runtime이다.

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

## Service Boundaries

- `../port-api` owns user authentication, durable conversation state, media room
  reservation, participant token issuance, and browser-facing event fanout.
- `../port-gateway` owns internal dispatch handling, service authentication,
  agent attach orchestration, recording orchestration, SIP orchestration, and
  lifecycle reporting back to API.
- `../port-media` owns SFU rooms, signaling, participants, tracks, WebRTC
  forwarding, and live media runtime state.
- `../port-record` owns recording finalize, encode, encrypt, retention, and
  storage workflows.
- `port-voice-pipeline` owns only the agent runtime: consume media as an agent
  participant, run STT -> LLM -> TTS, publish synthesized audio, and expose a
  narrow internal attach/control boundary for gateway.

## Gateway Integration

Gateway is the only service that should attach a voice pipeline to a reserved
conversation. The expected control flow is:

1. API creates a conversation, reserves a media room, and mints an internal
   agent participant token.
2. API dispatches the session to gateway through its internal dispatch endpoint.
3. Gateway validates service auth and asks `port-voice-pipeline` to attach the
   agent participant using `conversationId`, `sessionId`, `roomId`,
   `mediaSignalingUrl`, `participantId`, and `participantToken`.
4. Voice pipeline joins `port-media` using the provided participant token and
   starts the voice pipeline.
5. Gateway reports `agent.started` or `agent.failed` to API. Voice pipeline does
   not report durable browser-facing lifecycle events directly.

The pipeline-side attach API must remain internal-service only. It should not
accept browser credentials, user tokens, or durable conversation mutations.

## Responsibilities

### Voice Pipeline Owns

- Agent participant runtime after gateway attach.
- STT, LLM, TTS provider adapter wiring.
- Audio ingress/egress for the agent participant.
- Turn-taking, barge-in, endpointing, and provider retry behavior inside the
  voice session.
- Process/session cancellation when gateway or media closes the session.

### Voice Pipeline Does Not Own

- API dispatch intake.
- User authentication or authorization.
- Durable conversation or transcript storage.
- Browser SSE/WebSocket fanout.
- Recording storage or retention.
- SIP/PSTN call legs.
- Media room ownership or SFU forwarding.

## Shared Constraints

- Security:
  - Provider API key는 환경변수로만 읽는다.
  - domain/application 계층에 secret을 전달하지 않는다.
  - Gateway-provided participant tokens are runtime credentials and must not be
    logged.
  - Internal attach/control endpoints require service-to-service authentication.
- Reliability:
  - 모든 long-running loop는 `context.Context` 취소를 따른다.
  - provider 오류는 wrapping해서 세션 경계에서 처리한다.
  - Attach failure must be returned to gateway so gateway can emit
    `agent.failed`.
- Performance:
  - audio frame은 작은 immutable value로 다룬다.
  - MVP는 backpressure를 channel buffer와 context cancellation으로 제어한다.
- Operational limits:
  - Pion media adapter는 MVP에서 skeleton/mock 중심이다.
  - 실제 codec/resampling 최적화는 v2에서 확정한다.
  - Gateway integration may be implemented before production media track binding,
    but the boundary must stay compatible with gateway dispatch fields.
