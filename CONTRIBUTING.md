# Contributing to gpc

Thanks for your interest in contributing! This document covers how to get started.

## Requirements

- Go 1.23 or later
- [golangci-lint](https://golangci-lint.run/welcome/install/)
- GNU Make (or compatible)

## Setup

1. Fork and clone the repository:

```bash
git clone https://github.com/<your-username>/gpc.git
cd gpc
```

2. Install the git hooks:

```bash
git config core.hooksPath .githooks
```

3. Verify the build:

```bash
make build
make test
```

## Development Workflow

1. Create a branch from `main`:

```bash
git checkout -b feature/my-feature
```

2. Make your changes and ensure everything passes:

```bash
make check    # runs fmt, vet, lint, test
```

3. Commit your changes with a clear message.

4. Push and open a pull request against `main`.

## Code Standards

- Follow existing patterns in the codebase
- Use `gofmt -s` for formatting (enforced by CI)
- All exported types and functions need doc comments
- Errors should wrap context: `fmt.Errorf("could not do X: %w", err)`
- Use the `exitcode` package for exit codes - never call `os.Exit()` directly

## Makefile Targets

| Target | Description |
|---|---|
| `make build` | Build the binary with version info |
| `make test` | Run tests with race detection |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code |
| `make vet` | Run go vet |
| `make check` | Run all checks (fmt, vet, lint, test) |
| `make clean` | Remove build artifacts |

## Testing

- Write table-driven tests where possible
- Use `NewClientWithHTTP` for API tests with `httptest.Server`
- Test both JSON and table output paths
- Run tests with race detection: `go test -race ./...`

## Adding a New Command

1. Create a new package under `cmd/<command>/`
2. Implement `NewCmd() *cobra.Command` that returns the command tree
3. Register it in `cmd/root.go` under the appropriate section
4. Use `cmdutil.ResolvePackage()` for the package flag
5. Use `cmdutil.GetOutputFormat()` and `output.Print()` for output
6. Use `cmdutil.RequireAuth()` for authenticated commands
7. Return `exitcode` errors, not raw errors
8. Add path builder functions to `internal/api/client.go` if needed

## Questions?

Open an issue or start a discussion on the [GitHub Discussions](https://github.com/andresdefi/Google-Play-Console-CLI/discussions) page.
