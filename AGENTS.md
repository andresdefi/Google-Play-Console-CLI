# gpc - Google Play Console CLI

## Project Overview

- **Module**: `github.com/andresdefi/gpc`
- **Binary**: `gpc`
- **Repo**: `andresdefi/Google-Play-Console-CLI`
- **Go version**: 1.26+
- **CLI framework**: Cobra (`github.com/spf13/cobra`)
- **Status**: Active development, CI green, all tests passing

## Architecture

```
main.go              -> cmd.Execute() returns exit code
cmd/root.go          -> Root command, grouped help, global flags (--package, --output),
                        fuzzy suggestions, registers all subcommands
cmd/<resource>/      -> One package per API resource group, each exports NewCmd()
cmd/config/          -> CLI configuration management (set/get/list/path)
cmd/vitals/          -> App Vitals via Play Developer Reporting API
internal/api/        -> HTTP client, path builders, edit helpers, upload support
internal/auth/       -> Service account JSON + OAuth2 token exchange + keyring storage
internal/config/     -> ~/.gpc/config.json management (atomic writes, 0600 perms)
internal/output/     -> TTY detection, table/JSON/CSV/YAML output routing
internal/exitcode/   -> Exit code constants (0-5) and ExitError type
internal/cmdutil/    -> ResolvePackage(), GetOutputFormat(), RequireAuth()
internal/version/    -> Version/Commit/Date variables injected via ldflags
```

## Stats

- **57 Go source files** across 39 packages
- **213 test functions** across 9 test files
- **~130 API endpoints** covered across 36 command groups
- **4 output formats**: table, JSON, CSV, YAML

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
1. `cmdutil.ResolvePackage(cmd)` - get package name from --package flag or config
2. `cmdutil.RequireAuth()` - get OAuth2 token or fail with exit code 3
3. `api.NewClient(token)` - create API client
4. Do the work, return `exitcode.*Error()` on failure
5. `output.Print(format, data, tableRenderer)` - output results

### Edit Pattern
The Google Play API uses "edits" (transactions) for app changes. Two patterns:
- **Read-only**: Create temp edit, read data, delete edit (see `withTempEdit` in tracks/listings)
- **Write**: `client.WithEdit(pkg, func(editID) error)` - auto creates, commits on success, deletes on error

### Path Builders
All API URL paths are built by functions in `internal/api/client.go` (e.g. `api.TracksPath(pkg, editID)`). Never hardcode paths in commands.

### Grouped Help
`cmd/root.go` overrides Cobra's default help with `groupedHelp()` that organizes commands into 10 logical categories (Getting Started, Release Pipeline, Monetization, App Vitals, etc.) using `text/tabwriter`.

### Fuzzy Suggestions
`rootCmd.SuggestionsMinimumDistance = 2` enables "Did you mean...?" on command typos.

### Exit Codes
- 0: Success
- 1: General error
- 2: Usage error
- 3: Auth error
- 4: API error
- 5: Config error

### Output Formats
- `--output table` (default in TTY) - go-pretty tables
- `--output json` (default in pipes) - indented JSON
- `--output csv` - CSV with headers
- `--output yaml` - YAML via gopkg.in/yaml.v3
- Auto-detection via `term.IsTerminal()`

## Build & CI

```bash
make build      # Binary with version metadata via ldflags
make test       # Tests with race detection
make lint       # golangci-lint v2
make check      # fmt + vet + lint + test
make install    # Install to GOPATH/bin
```

### CI/CD (GitHub Actions)
- **ci.yml**: Build (Go 1.26 + stable matrix), lint (golangci-lint v2.11.4 via action v7)
- **release.yml**: GoReleaser on tag push, Homebrew tap auto-update
- **security.yml**: CodeQL + govulncheck, weekly schedule

### Linting
golangci-lint v2 with: errcheck, govet, ineffassign, staticcheck, unused, misspell

## Testing Patterns

- API client: `httptest.NewServer` with mock handlers, `NewClientWithHTTP` for injection
- Config: `t.TempDir()` + `t.Setenv("HOME", ...)` for filesystem isolation
- Commands: verify command tree structure via `rootCmd.Commands()`
- Output: capture stdout/stderr via `os.Pipe()`
- Auth: keychain bypassed in tests via config file fallback

## API Coverage

Full coverage of Google Play Developer API v3 (~130 endpoints across 36 resource groups):

**Release Pipeline**: edits, tracks, releases (with deploy convenience), apks, bundles,
deobfuscation, expansion-files, country-availability

**Monetization**: iap, subscriptions, base-plans, offers, one-time-products,
purchase-options, otp-offers, pricing

**Store Presence**: listings, images, details, testers, reviews, data-safety

**App Vitals** (Play Developer Reporting API): crashes, anrs, startup, rendering,
battery, errors (counts + issues)

**Orders & Purchases**: orders, purchases (products v1/v2, subscriptions v1/v2, voided)

**Account**: users, grants

**Device & Recovery**: devices, app-recovery, external-transactions

**APK Variants**: generated-apks, system-apks, internal-sharing

## Conventions

- Named exports preferred
- Errors wrap context: `fmt.Errorf("could not X: %w", err)`
- Commands return `exitcode.*Error()`, never call `os.Exit()` directly
- Use existing path builders in `internal/api/client.go`, don't hardcode API URLs
- Unchecked error returns: use `_ =` or `defer func() { _ = ... }()`
- Table renderers go inline in command RunE functions
- stdin JSON for complex request bodies (create/update commands)
- `google.CredentialsFromJSONWithType` (not deprecated `CredentialsFromJSON`)
