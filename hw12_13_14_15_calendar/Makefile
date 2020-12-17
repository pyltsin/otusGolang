CMD := calendar

build:
	go build -o $(CMD) ./cmd/calendar/main.go ./cmd/calendar/version.go

test:
	go test -race ./internal/...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...


.PHONY: build test lint

# BIN := "./bin/calendar"
# DOCKER_IMG="calendar:develop"
#
# GIT_HASH := $(shell git log --format="%h" -n 1)
# LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
#
# build:
# 	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
#
# run: build
# 	$(BIN) -config ./configs/config.toml
#
# build-img:
# 	docker build \
# 		--build-arg=LDFLAGS="$(LDFLAGS)" \
# 		-t $(DOCKER_IMG) \
# 		-f build/Dockerfile .
#
# run-img: build-img
# 	docker run $(DOCKER_IMG)
#
# version: build
# 	$(BIN) version
#
# test:
# 	go test -race ./internal/... ./pkg/...
#
# install-lint-deps:
# 	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.33.0
#
# lint: install-lint-deps
# 	golangci-lint run ./...
#
# .PHONY: build run build-img run-img version test lint
