# gpc - Google Play Console CLI
# https://github.com/andresdefi/Google-Play-Console-CLI

BINARY    := gpc
MODULE    := github.com/andresdefi/gpc
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo "devel")
COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE      := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS   := -s -w \
	-X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.Date=$(DATE)

GOLANGCI_LINT_VERSION := v2.11.4
GOFUMPT_VERSION       := v0.9.2
GOSEC_VERSION         := v2.22.3

.PHONY: build install test test-coverage lint fmt vet check security tools \
        generate-docs install-hooks clean help

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

## test-integration: Run integration tests (requires GPC_KEY_FILE env var)
test-integration:
	GPC_INTEGRATION_TEST=1 go test -race -run Integration ./...

## lint: Run golangci-lint
lint:
	golangci-lint run

## fmt: Format all Go source files with gofumpt
fmt:
	gofumpt -l -w .

## vet: Run go vet
vet:
	go vet ./...

## security: Run gosec security scanner
security:
	gosec -quiet ./...

## check: Run fmt, vet, lint, security, and test (sequential)
check: fmt vet lint security test

## tools: Install development tool dependencies
tools:
	go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION)

## generate-docs: Generate command reference documentation
generate-docs: build
	@mkdir -p docs
	@echo "# gpc Command Reference" > docs/COMMANDS.md
	@echo "" >> docs/COMMANDS.md
	@echo "Auto-generated from \`gpc --help\` on $$(date -u +%Y-%m-%d)." >> docs/COMMANDS.md
	@echo "" >> docs/COMMANDS.md
	@./$(BINARY) --help >> docs/COMMANDS.md 2>&1
	@echo "" >> docs/COMMANDS.md
	@for cmd in $$(./$(BINARY) --help 2>&1 | grep -E '^\s+[a-z]' | awk '{print $$1}' | sort -u | grep -v gpc); do \
		echo "## $$cmd" >> docs/COMMANDS.md; \
		echo "" >> docs/COMMANDS.md; \
		echo '```' >> docs/COMMANDS.md; \
		./$(BINARY) $$cmd --help >> docs/COMMANDS.md 2>&1 || true; \
		echo '```' >> docs/COMMANDS.md; \
		echo "" >> docs/COMMANDS.md; \
	done
	@echo "Generated docs/COMMANDS.md"

## install-hooks: Configure git pre-commit hooks
install-hooks:
	git config core.hooksPath .githooks

## clean: Remove build artifacts and coverage files
clean:
	rm -f $(BINARY) coverage.out coverage.html

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
