# State

current_version: v2

## Active

- Version: v2
- Focus: voice pipeline stability and hexagonal wiring cleanup
- Current feature: pipeline-stability-hex-wiring

## Completed Versions

### v1

- Pion RTP 기반 voice agent worker MVP skeleton을 DDD + Hexagonal 경계로 구성했다.
- `internal/domain/voice`에 PCM frame, transcript, assistant response, speech activity event value를 정의했다.
- `internal/application/ports`에 audio ingress/egress, STT, LLM, TTS, VAD 포트를 분리했다.
- `internal/application/session`에 STT -> LLM -> TTS -> audio egress orchestration, 반복 실행 runner, turn-aware orchestration을 추가했다.
- `internal/application/turn`에 barge-in과 endpointing decision controller/activity processor를 추가했다.
- `internal/adapters/pionrtp`와 `internal/adapters/media/pion`에 Pion RTP ingress/egress, Opus decode boundary, RTP packetization, runner binder를 구성했다.
- `internal/adapters/providers`에 Deepgram STT, OpenRouter LLM, Cartesia TTS adapter와 provider factory wiring을 추가했다.
- `internal/adapters/vad`에 noop/Silero engine 기반 VAD adapter와 `VAD_PROVIDER` wiring을 추가했다.
- `cmd/worker`에서 config loading, provider wiring, turn runtime 선택, `RUN_SESSION` runner 실행 경로를 연결했다.
- 검증: `go test ./...`, `make build`.
