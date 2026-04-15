# gpc - Google Play Console CLI

[![Go Report Card](https://goreportcard.com/badge/github.com/andresdefi/gpc)](https://goreportcard.com/report/github.com/andresdefi/gpc)
[![CI](https://github.com/andresdefi/Google-Play-Console-CLI/actions/workflows/ci.yml/badge.svg)](https://github.com/andresdefi/Google-Play-Console-CLI/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/andresdefi/Google-Play-Console-CLI)](https://github.com/andresdefi/Google-Play-Console-CLI/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/andresdefi/gpc.svg)](https://pkg.go.dev/github.com/andresdefi/gpc)

Fast, lightweight, scriptable CLI for the Google Play Developer API.

`gpc` gives you complete coverage of the Android Publisher API v3 from your terminal - manage apps, releases, in-app products, subscriptions, reviews, vitals, and more. Pipe output into `jq`, use it in CI/CD pipelines, or just skip the web console.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Output Formatting](#output-formatting)
- [The Deploy Flow](#the-deploy-flow)
- [Configuration](#configuration)
- [Environment Variables](#environment-variables)
- [Diagnostics](#diagnostics)
- [Exit Codes](#exit-codes)
- [Shell Completions](#shell-completions)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Homebrew (macOS/Linux)

```bash
brew install andresdefi/tap/gpc
```

### Install script

```bash
curl -sSfL https://raw.githubusercontent.com/andresdefi/Google-Play-Console-CLI/main/install.sh | bash
```

### Go install

```bash
go install github.com/andresdefi/gpc@latest
```

### Binary download

Download the latest release from the [releases page](https://github.com/andresdefi/Google-Play-Console-CLI/releases/latest) and add it to your `PATH`.

## Quick Start

```bash
# Authenticate with a service account key
gpc auth login --key-file credentials.json

# Set a default package to avoid repeating -p
gpc config set package com.example.app

# Check your setup
gpc doctor

# Deploy a bundle to the internal track
gpc releases deploy app-release.aab --track internal

# Promote from beta to production
gpc releases promote --from beta --to production

# Check crash rates
gpc vitals crashes

# List reviews
gpc reviews list
```

## Command Reference

Commands are organized into logical groups. Run `gpc --help` to see them all, or see the full [command reference](docs/COMMANDS.md).

| Group | Commands | Description |
|-------|----------|-------------|
| **Getting Started** | `auth`, `config`, `doctor`, `version` | Authentication, configuration, diagnostics |
| **App Management** | `apps`, `edits` | App details and edit sessions |
| **Release Pipeline** | `releases`, `tracks`, `apks`, `bundles`, `deobfuscation`, `expansionfiles`, `countryavailability` | Build, upload, and release management |
| **Monetization** | `iap`, `subscriptions`, `baseplans`, `offers`, `onetimeproducts`, `purchaseoptions`, `otpoffers`, `pricing` | In-app products, subscriptions, pricing |
| **Store Presence** | `listings`, `images`, `details`, `testers`, `reviews`, `datasafety` | Store listings, screenshots, reviews |
| **App Vitals** | `vitals` | Crash rates, ANR rates, startup, rendering, battery |
| **Orders & Purchases** | `orders`, `purchases` | Order management, purchase verification |
| **Account Management** | `users`, `grants` | Developer account users and permissions |
| **Device & Recovery** | `devices`, `apprecovery`, `externaltransactions` | Device tiers, recovery actions |
| **APK Variants** | `generatedapks`, `systemapks`, `internalsharing` | Generated APKs, system images, internal sharing |

## Output Formatting

gpc automatically detects whether stdout is a terminal:

- **Terminal (TTY)**: renders human-friendly tables
- **Pipe/redirect**: outputs JSON for scripting

Four output formats are supported:

```bash
gpc reviews list --output table   # Human-readable tables (default in terminal)
gpc reviews list --output json    # JSON (default when piped)
gpc reviews list --output csv     # CSV with headers
gpc reviews list --output yaml    # YAML
```

Supports `NO_COLOR` environment variable to disable table styling ([no-color.org](https://no-color.org/)).

## The Deploy Flow

The `releases deploy` command wraps the full Google Play edit flow into a single step with progress feedback:

1. Creates an edit session
2. Uploads your APK or AAB (with spinner)
3. Assigns the artifact to the specified track
4. Commits the edit

```bash
# Full rollout to production
gpc releases deploy app-release.aab --track production

# Staged rollout to 10% of users
gpc releases deploy app.aab --track production --rollout 0.1

# With release name and notes
gpc releases deploy app.aab --track beta \
  --release-name "v2.0.0" --notes "New features and bug fixes"
```

For more control, use the lower-level commands: `edits`, `bundles`/`apks`, and `tracks`.

## Configuration

gpc stores configuration in `~/.gpc/config.json`. Manage it with the `config` command:

```bash
gpc config set package com.example.app   # Set default package
gpc config get package                    # Get a value
gpc config list                           # Show all settings
gpc config path                           # Print config file path
```

| Setting | Description |
|---------|-------------|
| `key_file_path` | Path to service account JSON key file |
| `package_name` | Default Android package name |

## Environment Variables

Environment variables override config file values. Flags override both.

| Variable | Description |
|----------|-------------|
| `GPC_PACKAGE` | Default package name (overrides config) |
| `GPC_KEY_FILE` | Path to service account key file |
| `GPC_OUTPUT` | Default output format (`json`, `table`, `csv`, `yaml`) |
| `NO_COLOR` | Disable colored/styled output ([no-color.org](https://no-color.org/)) |

Priority: `--flag` > `env var` > `config file` > `auto-detect`

## Diagnostics

Run `gpc doctor` to validate your setup:

```bash
$ gpc doctor
Running diagnostics...

  + Version: gpc v0.1.0 (commit: abc1234, built: 2026-04-15T00:00:00Z)
  + Go runtime: go1.26.2 darwin/arm64
  + Config file: /Users/you/.gpc/config.json
  + Config readable: default package: com.example.app
  + Credentials: mysa...vice@project.iam.gserviceaccount.com
  + OAuth2 token: valid
  + API reachable: ok (142ms)
  + Environment: no env overrides set

8 passed, 0 failed
```

## Exit Codes

gpc uses granular exit codes for scripting:

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Usage error (invalid arguments) |
| 3 | Authentication error (401/403) |
| 4 | Not found (404) |
| 5 | Conflict (409) |
| 6 | Configuration error |
| 10-59 | HTTP 4xx errors (code = 10 + status - 400) |
| 60-99 | HTTP 5xx errors (code = 60 + status - 500) |

## Shell Completions

Generate shell completions for your shell:

```bash
# Bash
gpc completion bash > /etc/bash_completion.d/gpc

# Zsh
gpc completion zsh > "${fpath[1]}/_gpc"

# Fish
gpc completion fish > ~/.config/fish/completions/gpc.fish

# PowerShell
gpc completion powershell > gpc.ps1
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, guidelines, and how to submit changes.

## License

[MIT](LICENSE)

---

**Disclaimer**: gpc is not affiliated with, endorsed by, or sponsored by Google. Google Play and the Google Play logo are trademarks of Google LLC.
