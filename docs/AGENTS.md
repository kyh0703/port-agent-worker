# Project Instructions

- 사용자에게 보여주는 설명은 한글로 작성한다.
- 구현은 DDD + Hexagonal Architecture 경계를 유지한다.
- domain/application 포트는 외부 SDK 타입을 노출하지 않는다.
- Pion, Deepgram, OpenRouter, TTS provider 연동은 adapter 계층에 둔다.
- MVP에서는 녹취, RAG, filler, smart turn analyzer 구현을 제외하고 인터페이스 확장 지점만 둔다.
- 실제 구현은 `.worktrees/` linked worktree에서 수행한다.

