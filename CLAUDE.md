# gpc - Google Play Console CLI

## Project Overview

- **Module**: `github.com/andresdefi/gpc`
- **Binary**: `gpc`
- **Repo**: `andresdefi/Google-Play-Console-CLI`
- **Go version**: 1.23+
- **CLI framework**: Cobra (`github.com/spf13/cobra`)

## Architecture

```
main.go              -> cmd.Execute() returns exit code
cmd/root.go          -> Root command, global flags (--package, --output), registers all subcommands
cmd/<resource>/      -> One package per API resource group, each exports NewCmd()
internal/api/        -> HTTP client, path builders, edit helpers, upload support
internal/auth/       -> Service account JSON + OAuth2 token exchange + keyring storage
internal/config/     -> ~/.gpc/config.json management
internal/output/     -> TTY detection, table (go-pretty) / JSON output routing
internal/exitcode/   -> Exit code constants (0-5) and ExitError type
internal/cmdutil/    -> ResolvePackage(), GetOutputFormat(), RequireAuth()
internal/version/    -> Version/Commit/Date variables injected via ldflags
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

### Exit Codes
- 0: Success
- 1: General error
- 2: Usage error
- 3: Auth error
- 4: API error
- 5: Config error

### Output
- TTY detected: table format (go-pretty)
- Piped/redirected: JSON format
- Override: `--output json` or `--output table`

## Build

```bash
make build      # Binary with version metadata
make test       # Tests with race detection
make lint       # golangci-lint
make check      # fmt + vet + lint + test
```

Version injected via ldflags: `-X .../version.Version=... -X .../version.Commit=... -X .../version.Date=...`

## Testing

- API client: `httptest.NewServer` with mock handlers
- Config: `t.TempDir()` + `t.Setenv("HOME", ...)` for isolation
- Commands: verify command tree structure via `rootCmd.Commands()`
- Output: capture stdout/stderr via `os.Pipe()`

## API Coverage

Full coverage of Google Play Developer API v3 (~130 endpoints across 35 resource groups):
edits, tracks, releases, apks, bundles, listings, images, details, testers,
iap, subscriptions, base-plans, offers, one-time-products, purchase-options,
otp-offers, pricing, reviews, orders, purchases, users, grants, devices,
app-recovery, external-transactions, generated-apks, system-apks,
internal-sharing, data-safety, deobfuscation, expansion-files, country-availability

## Conventions

- Named exports preferred
- Errors wrap context: `fmt.Errorf("could not X: %w", err)`
- Commands return `exitcode.*Error()`, never call `os.Exit()` directly
- Use existing path builders, don't hardcode API URLs
- Table renderers go inline in command RunE functions
- stdin JSON for complex request bodies (create/update commands)
