CMD := calendar

generate:
	go generate ./...

build:
	go build -o $(CMD) ./cmd/calendar/main.go ./cmd/calendar/version.go

test:
	go test -race ./internal/...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...


.PHONY: build test lint generate