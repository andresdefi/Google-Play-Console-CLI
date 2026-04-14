package cmdutil

import (
	"fmt"

	"github.com/andresdefi/gpc/internal/auth"
	"github.com/andresdefi/gpc/internal/config"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

// ResolvePackage returns the package name from the flag or config default.
func ResolvePackage(cmd *cobra.Command) (string, error) {
	pkg, _ := cmd.Flags().GetString("package")
	if pkg != "" {
		return pkg, nil
	}

	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("could not load config: %w", err)
	}

	if cfg.PackageName != "" {
		return cfg.PackageName, nil
	}

	return "", fmt.Errorf("no package name specified; use --package flag or set a default with 'gpc config set package <name>'")
}

// GetOutputFormat returns the output format from the flag or auto-detection.
func GetOutputFormat(cmd *cobra.Command) output.Format {
	explicit, _ := cmd.Flags().GetString("output")
	return output.DetectFormat(explicit)
}

// RequireAuth loads the auth token and returns an authenticated API client token.
// Returns an error if not authenticated.
func RequireAuth() (string, error) {
	token, err := auth.GetToken()
	if err != nil {
		return "", fmt.Errorf("not authenticated; run 'gpc auth login' first: %w", err)
	}
	if token == "" {
		return "", fmt.Errorf("not authenticated; run 'gpc auth login' first")
	}
	return token, nil
}
