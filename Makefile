APP := port-voice-pipeline

.PHONY: test build tidy

test:
	go test ./...

build:
	go build -o bin/$(APP) ./cmd/worker

tidy:
	go mod tidy

