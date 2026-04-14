package cmdutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/andresdefi/gpc/internal/auth"
	"github.com/andresdefi/gpc/internal/config"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

// Environment variable names for configuration override.
// Priority: flag > env var > config file.
const (
	EnvPackage = "GPC_PACKAGE"
	EnvKeyFile = "GPC_KEY_FILE"
	EnvOutput  = "GPC_OUTPUT"
)

// ResolvePackage returns the package name from flag, env var, or config default.
func ResolvePackage(cmd *cobra.Command) (string, error) {
	// 1. Flag (highest priority).
	pkg, _ := cmd.Flags().GetString("package")
	pkg = strings.TrimSpace(pkg)
	if pkg != "" {
		return pkg, nil
	}

	// 2. Environment variable.
	if env := strings.TrimSpace(os.Getenv(EnvPackage)); env != "" {
		return env, nil
	}

	// 3. Config file.
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("could not load config: %w", err)
	}

	if cfg.PackageName != "" {
		return strings.TrimSpace(cfg.PackageName), nil
	}

	return "", fmt.Errorf("no package name specified; use --package flag, set GPC_PACKAGE env var, or run 'gpc config set package <name>'")
}

// GetOutputFormat returns the output format from flag, env var, or auto-detection.
func GetOutputFormat(cmd *cobra.Command) output.Format {
	// 1. Flag.
	explicit, _ := cmd.Flags().GetString("output")
	explicit = strings.TrimSpace(explicit)
	if explicit != "" {
		return output.DetectFormat(explicit)
	}

	// 2. Environment variable.
	if env := strings.TrimSpace(os.Getenv(EnvOutput)); env != "" {
		return output.DetectFormat(env)
	}

	// 3. Auto-detect.
	return output.DetectFormat("")
}

// RequireAuth loads the auth token and returns a validated OAuth2 access token.
// Returns an error if not authenticated or if the token is invalid.
func RequireAuth() (string, error) {
	token, err := auth.GetToken()
	if err != nil {
		return "", fmt.Errorf("not authenticated; run 'gpc auth login' first: %w", err)
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", fmt.Errorf("not authenticated; run 'gpc auth login' first")
	}
	return token, nil
}

// SanitizeArg trims whitespace from a user-provided argument and validates it's non-empty.
func SanitizeArg(value, name string) (string, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	return v, nil
}
