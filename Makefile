# gpc - Google Play Console CLI
# https://github.com/andresdefi/gpc

BINARY    := gpc
MODULE    := github.com/andresdefi/gpc
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "devel")
COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE      := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS   := -s -w \
	-X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.Date=$(DATE)

.PHONY: build install test test-coverage lint fmt vet check clean help

## build: Build the gpc binary
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

## install: Install gpc to $GOPATH/bin
install:
	go install -ldflags "$(LDFLAGS)" .

## test: Run all tests with race detection
test:
	go test -race -coverprofile=coverage.out ./...

## test-coverage: Open test coverage report in browser
test-coverage: test
	go tool cover -html=coverage.out

## lint: Run golangci-lint
lint:
	golangci-lint run

## fmt: Format all Go source files
fmt:
	gofmt -s -w .

## vet: Run go vet
vet:
	go vet ./...

## check: Run fmt, vet, lint, and test (sequential)
check: fmt vet lint test

## clean: Remove build artifacts and coverage files
clean:
	rm -f $(BINARY) coverage.out coverage.html

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
