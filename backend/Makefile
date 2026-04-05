.PHONY: dev test lint

dev:
	go run ./cmd/server

test:
	go test ./... -v -count=1

lint:
	go vet ./...
