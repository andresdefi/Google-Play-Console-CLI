# gpc - Google Play Console CLI

[![Go Report Card](https://goreportcard.com/badge/github.com/andresdefi/gpc)](https://goreportcard.com/report/github.com/andresdefi/gpc)
[![CI](https://github.com/andresdefi/Google-Play-Console-CLI/actions/workflows/ci.yml/badge.svg)](https://github.com/andresdefi/Google-Play-Console-CLI/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/andresdefi/Google-Play-Console-CLI)](https://github.com/andresdefi/Google-Play-Console-CLI/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/andresdefi/gpc.svg)](https://pkg.go.dev/github.com/andresdefi/gpc)

Fast, lightweight, scriptable CLI for the Google Play Developer API.

`gpc` gives you complete coverage of the Android Publisher API v3 from your terminal - manage apps, releases, in-app products, subscriptions, reviews, and more. Pipe JSON output into `jq`, use it in CI/CD pipelines, or just skip the web console.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Output Formatting](#output-formatting)
- [The Deploy Flow](#the-deploy-flow)
- [Configuration](#configuration)
- [Exit Codes](#exit-codes)
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

Download the latest release from the [releases page](https://github.com/andresdefi/gpc/releases/latest) and add it to your `PATH`.

## Quick Start

```bash
# Authenticate with a service account key
gpc auth login --key-file credentials.json

# List your apps
gpc apps get -p com.example.app

# Deploy a bundle to the internal track
gpc releases deploy app-release.aab -p com.example.app --track internal

# Promote from beta to production
gpc releases promote -p com.example.app --from beta --to production
```

## Command Reference

| Command | Description |
|---|---|
| `auth` | Authenticate with Google Play Developer API |
| `version` | Print the version of gpc |
| `apps` | Get application details |
| `edits` | Manage edit sessions (create, validate, commit, delete) |
| `tracks` | List and update release tracks |
| `releases` | Deploy, promote, rollout, and halt releases |
| `apks` | Upload and list APKs |
| `bundles` | Upload and list Android App Bundles |
| `deobfuscation` | Upload deobfuscation (mapping) files |
| `expansionfiles` | Manage APK expansion files |
| `countryavailability` | View track country availability |
| `iap` | Manage in-app products (one-time purchases) |
| `subscriptions` | Manage subscriptions |
| `baseplans` | Manage subscription base plans |
| `offers` | Manage subscription offers |
| `onetimeproducts` | Manage one-time products |
| `purchaseoptions` | Manage purchase options for one-time products |
| `otpoffers` | Manage offers for one-time products |
| `pricing` | Convert regional prices |
| `listings` | Manage store listings |
| `images` | Manage store listing images and screenshots |
| `details` | Manage app details (contact info, category) |
| `testers` | Manage track testers |
| `reviews` | List and reply to user reviews |
| `datasafety` | Manage data safety declarations |
| `orders` | Refund orders |
| `purchases` | Verify product and subscription purchases |
| `users` | Manage developer account users |
| `grants` | Manage user grants and permissions |
| `devices` | Manage device tier configurations |
| `apprecovery` | Manage app recovery actions |
| `externaltransactions` | Manage external transactions |
| `generatedapks` | Download generated APKs from bundles |
| `systemapks` | Manage system APK variants |
| `internalsharing` | Upload artifacts for internal app sharing |

## Output Formatting

gpc automatically detects whether stdout is a terminal:

- **Terminal (TTY)**: renders human-friendly tables
- **Pipe/redirect**: outputs JSON for scripting

Override with the `--output` flag:

```bash
# Force JSON output in a terminal
gpc apps get -p com.example.app --output json

# Force table output in a pipe
gpc apps get -p com.example.app --output table | less
```

## The Deploy Flow

The `releases deploy` command is a convenience that wraps the full Google Play edit flow into a single step:

1. Creates an edit session
2. Uploads your APK or AAB
3. Assigns the artifact to the specified track
4. Commits the edit

```bash
# Full rollout to production
gpc releases deploy app-release.aab -p com.example.app --track production

# Staged rollout to 10% of users
gpc releases deploy app.aab -p com.example.app --track production --rollout 0.1

# With release name and notes
gpc releases deploy app.aab -p com.example.app --track beta \
  --release-name "v2.0.0" --notes "New features and bug fixes"
```

For more control, use the lower-level commands: `edits`, `bundles`/`apks`, and `tracks`.

## Configuration

gpc stores configuration in `~/.gpc/config.json`.

```bash
# Set a default package name to avoid repeating -p
gpc config set package com.example.app

# Set the service account key file path
gpc auth login --key-file /path/to/credentials.json
```

| Setting | Description |
|---|---|
| `key_file_path` | Path to service account JSON key file |
| `package_name` | Default Android package name |

## Exit Codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | General error |
| 2 | Usage error (invalid arguments) |
| 3 | Authentication error |
| 4 | API error |
| 5 | Configuration error |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, guidelines, and how to submit changes.

## License

[MIT](LICENSE)

---

**Disclaimer**: gpc is not affiliated with, endorsed by, or sponsored by Google. Google Play and the Google Play logo are trademarks of Google LLC.
