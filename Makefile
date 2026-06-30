APP := port-agent-worker

.PHONY: test build tidy

test:
	go test ./...

build:
	go build -o bin/$(APP) ./cmd/worker

tidy:
	go mod tidy

