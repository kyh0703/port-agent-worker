# Roadmap

## Current Track

- Active version: v1
- Goal: Pion SFU audio track을 받아 STT -> LLM -> TTS -> RTP publish 흐름을 구성할 수 있는 voice pipeline skeleton을 만든다.
- Exit criteria:
  - Go service가 빌드되고 테스트된다.
  - RTP/PCM media boundary가 포트/어댑터로 분리된다.
  - STT, LLM, TTS provider가 교체 가능한 포트로 정의된다.
  - 세션 오케스트레이터가 DDD/Hexagonal 경계 안에서 조립된다.

## Future Versions

- `v2`:
  - Goal: Deepgram/OpenRouter/TTS provider 실연동, barge-in, endpointing 강화.
  - Dependencies: v1 voice pipeline skeleton.
- `v3`:
  - Goal: RAG, filler, emotion tag, metrics, recording.
  - Dependencies: v2 streaming flow.

## Deferred

- Silero VAD ONNX runtime integration.
- LocalSmartTurnAnalyzerV3.
- GStreamer fallback media adapter.
- Recording pipeline.
