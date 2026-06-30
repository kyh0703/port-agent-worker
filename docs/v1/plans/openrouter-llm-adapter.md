# OpenRouter LLM Adapter

## Goal

OpenRouter chat completions adapter를 추가해서 user utterance를 assistant response text로 변환한다.

## References

- `docs/STATE.md`
- `docs/ROADMAP.md`
- `docs/ARCHITECTURE.md`
- `docs/v1/designs/2026-06-30-v1-openrouter-llm-adapter.md`
- OpenRouter chat completions docs.

## Workspace

- Branch: feature/v1-openrouter-llm-adapter
- Base: main
- Isolation: `.worktrees/feature-v1-openrouter-llm-adapter`
- Created by: exec-plan via git-worktree

## Task Graph

### Task T1

- [ ] Complete
- Goal: OpenRouter adapter config와 chat completion request/response 모델을 구현한다.
- Depends on:
  - none
- Write Scope:
  - `internal/adapters/providers/openrouter/**`
- Read Context:
  - `internal/application/ports/providers.go`
  - `internal/domain/voice/**`
- Checks:
  - `go test ./...`
- Parallel-safe: no

### Task T2

- [ ] Complete
- Goal: HTTP call, response parsing, error handling 테스트를 구현한다.
- Depends on:
  - T1
- Write Scope:
  - `internal/adapters/providers/openrouter/**`
- Read Context:
  - `internal/domain/voice/conversation.go`
- Checks:
  - `go test ./internal/adapters/providers/openrouter`
- Parallel-safe: no

### Task T3

- [ ] Complete
- Goal: 전체 빌드 검증을 완료한다.
- Depends on:
  - T1
  - T2
- Write Scope:
  - `internal/adapters/providers/openrouter/**`
- Read Context:
  - full tree
- Checks:
  - `go test ./...`
  - `make build`
- Parallel-safe: no

