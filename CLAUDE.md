# gpc - Google Play Console CLI

## Project Overview

- **Module**: `github.com/andresdefi/gpc`
- **Binary**: `gpc`
- **Repo**: `andresdefi/Google-Play-Console-CLI`
- **Go version**: 1.26+
- **CLI framework**: Cobra (`github.com/spf13/cobra`)
- **Formatter**: gofumpt (not gofmt)
- **Latest release**: v0.1.0
- **Status**: v0.1.0 released, CI green, 244 tests passing

## Architecture

```
main.go              -> cmd.Execute() returns exit code
cmd/root.go          -> Root command, grouped help (10 categories), global flags
                        (--package, --output), fuzzy suggestions, registers all subcommands
cmd/<resource>/      -> One package per API resource group, each exports NewCmd()
cmd/config/          -> CLI configuration management (set/get/list/path)
cmd/doctor/          -> Diagnostics (config, auth, API reachability, env vars)
cmd/vitals/          -> [beta] App Vitals via Play Developer Reporting API
internal/api/        -> HTTP client, path builders, edit helpers, upload, pagination
internal/api/pagination.go -> ListAll/ListAllRaw with nextPageToken handling
internal/auth/       -> Service account JSON + OAuth2 token exchange + keyring storage
internal/config/     -> ~/.gpc/config.json management (atomic writes, 0600 perms)
internal/output/     -> TTY detection, table/JSON/CSV/YAML output, NO_COLOR support
internal/exitcode/   -> Granular exit codes (0-6 + HTTP ranges 10-59, 60-99)
internal/cmdutil/    -> ResolvePackage(), GetOutputFormat(), RequireAuth(), SanitizeArg()
internal/spinner/    -> TTY-aware progress spinner for long operations
internal/testutil/   -> Integration test helpers (SkipUnlessIntegration, RequireKeyFile)
internal/version/    -> Version/Commit/Date variables injected via ldflags
docs/COMMANDS.md     -> Auto-generated command reference (make generate-docs)
```

## Stats

- **64 Go source files** across 40+ packages
- **244 test functions** across 12 test files
- **~130 API endpoints** covered across 37 command groups
- **4 output formats**: table, JSON, CSV, YAML
- **5 release binaries**: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64

## Dependencies

```
github.com/spf13/cobra           # CLI framework
github.com/jedib0t/go-pretty/v6  # Table output
github.com/zalando/go-keyring    # Credential storage (keychain)
golang.org/x/oauth2              # Google OAuth2 token exchange
golang.org/x/term                # TTY detection
gopkg.in/yaml.v3                 # YAML output
```

## Key Patterns

### Command Pattern
Every command package exports `NewCmd() *cobra.Command`. Inside:
1. `cmdutil.ResolvePackage(cmd)` - get package from flag > env > config
2. `cmdutil.RequireAuth()` - get OAuth2 token, validate non-empty, or fail with exit code 3
3. `api.NewClient(token)` - create API client
4. Do the work, return `exitcode.*Error()` on failure
5. `output.Print(format, data, tableRenderer)` - output results

### Edit Pattern
The Google Play API uses "edits" (transactions) for app changes. Two patterns:
- **Read-only**: Create temp edit, read data, delete edit (see `withTempEdit` in tracks/listings)
- **Write**: `client.WithEdit(pkg, func(editID) error)` - auto creates, commits on success, deletes on error

### Path Builders
All API URL paths are built by functions in `internal/api/client.go`. Never hardcode paths in commands.

### Pagination
`client.ListAll(path, params, mergeFn)` follows `nextPageToken` across all pages. Use for any list endpoint.

### Config Resolution
Priority: `--flag` > env var (`GPC_PACKAGE`, `GPC_KEY_FILE`, `GPC_OUTPUT`) > `~/.gpc/config.json` > auto-detect

### Exit Codes
- 0: Success, 1: General error, 2: Usage error
- 3: Auth error (401/403), 4: Not found (404), 5: Conflict (409), 6: Config error
- 10-59: HTTP 4xx (code = 10 + status - 400)
- 60-99: HTTP 5xx (code = 60 + status - 500)

### Output Formats
- `--output table` (default in TTY) - go-pretty tables, respects NO_COLOR
- `--output json` (default in pipes) - indented JSON
- `--output csv` - CSV with headers via `output.PrintWithCSV()`
- `--output yaml` - YAML via gopkg.in/yaml.v3
- Auto-detection via `term.IsTerminal()`

### Retry Logic
Retries on 429/500/502/503/504 with exponential backoff. Respects `Retry-After` header when present, capped at 60s.

### Stability Labels
Commands may be annotated with `[beta]` in their Short description to indicate pre-stable APIs (e.g. vitals).

## Build & CI

```bash
make build           # Binary with version metadata via ldflags
make test            # Tests with race detection
make test-integration # Integration tests (requires GPC_KEY_FILE + GPC_TEST_PACKAGE)
make lint            # golangci-lint v2
make fmt             # gofumpt formatting
make security        # gosec security scanner
make check           # fmt + vet + lint + security + test
make tools           # Install dev dependencies (golangci-lint, gofumpt, gosec)
make generate-docs   # Auto-generate docs/COMMANDS.md from CLI help
make install-hooks   # Configure git pre-commit hooks
```

### CI/CD (GitHub Actions)
- **ci.yml**: Build (Go 1.26 + stable matrix), lint (golangci-lint v2.11.4 via action v7)
- **release.yml**: GoReleaser on tag push (v*), cross-platform binaries + checksums
- **security.yml**: CodeQL + govulncheck, weekly schedule
- **Homebrew tap**: `andresdefi/homebrew-tap` - commented out in .goreleaser.yaml until HOMEBREW_TAP_TOKEN secret is configured

### Release Process
1. Ensure CI is green on main
2. `git tag v0.x.x && git push origin v0.x.x`
3. GoReleaser builds binaries for 5 platforms, creates GitHub Release with changelog

## Testing Patterns

- API client: `httptest.NewServer` with mock handlers, `NewClientWithHTTP` for injection
- Pagination: mock multi-page responses with nextPageToken
- Retry-After: mock 429 responses with header, verify backoff behavior
- Config: `t.TempDir()` + `t.Setenv("HOME", ...)` for filesystem isolation
- Env vars: `t.Setenv("GPC_PACKAGE", ...)` for env var override tests
- Commands: verify command tree structure, grouped help, fuzzy suggestions
- Output: capture stdout/stderr via `os.Pipe()`
- Integration: `testutil.SkipUnlessIntegration(t)` for opt-in real API tests
  (set `GPC_INTEGRATION_TEST=1`, `GPC_KEY_FILE`, `GPC_TEST_PACKAGE`)

## Conventions

- Format with gofumpt, not gofmt
- Named exports preferred
- Errors wrap context: `fmt.Errorf("could not X: %w", err)`
- Commands return `exitcode.*Error()`, never call `os.Exit()` directly
- Use existing path builders, don't hardcode API URLs
- Unchecked error returns: use `_ =` or `defer func() { _ = ... }()`
- Input sanitization: `cmdutil.SanitizeArg()` or `strings.TrimSpace()` on all user input
- `google.CredentialsFromJSONWithType` (not deprecated `CredentialsFromJSON`)
- Use `spinner.New()` for long operations (uploads, API calls)
- Don't delete AGENTS.md - it serves as context for other AI agents

## TODO

- [ ] Configure HOMEBREW_TAP_TOKEN secret and re-enable brew section in .goreleaser.yaml
- [ ] Write real integration tests once test app is set up in Play Console
- [ ] Consider a docs site (Mintlify or GitHub Pages) for workflow guides
