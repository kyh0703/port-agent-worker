# port-agent-worker

Pion SFU 기반 RTP audio track을 받아 STT -> LLM -> TTS 흐름을 실행하는 voice agent worker입니다.

## Architecture

- DDD + Hexagonal Architecture.
- Domain/application 계층은 Pion/provider SDK 타입을 모릅니다.
- Pion RTP, STT, LLM, TTS는 adapter로 연결합니다.

## MVP

```text
Pion RTP ingress
 -> PCM frame
 -> STT port
 -> LLM port
 -> TTS port
 -> PCM egress
 -> Pion RTP publish
```

## Commands

```bash
make test
make build
```

