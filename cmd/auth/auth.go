package auth

import (
	"fmt"
	"path/filepath"

	internalauth "github.com/andresdefi/gpc/internal/auth"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

// NewCmd returns the auth command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
		Long:  "Authenticate with the Google Play Developer API using a service account key file.",
	}

	cmd.AddCommand(newLoginCmd())
	cmd.AddCommand(newStatusCmd())
	cmd.AddCommand(newLogoutCmd())
	return cmd
}

func newLoginCmd() *cobra.Command {
	var keyFile string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a service account key file",
		Long:  "Store a Google Cloud service account key file for API authentication.",
		Example: `  gpc auth login --key-file service-account.json
  gpc auth login --key-file ~/keys/play-api.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if keyFile == "" {
				return exitcode.AuthError("--key-file is required")
			}

			absPath, err := filepath.Abs(keyFile)
			if err != nil {
				return exitcode.AuthError("could not resolve key file path: %v", err)
			}

			info, err := internalauth.Login(absPath)
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			output.Success(fmt.Sprintf("Authenticated as %s (project: %s)", info.ClientEmail, info.ProjectID))
			return nil
		},
	}

	cmd.Flags().StringVar(&keyFile, "key-file", "", "Path to service account JSON key file (required)")
	return cmd
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := internalauth.GetStatus()
			if err != nil {
				return exitcode.AuthError("not authenticated; run 'gpc auth login' first")
			}

			fmt.Printf("Account:  %s\n", internalauth.MaskEmail(info.ClientEmail))
			fmt.Printf("Project:  %s\n", info.ProjectID)

			// Verify token is still valid.
			if _, err := internalauth.GetToken(); err != nil {
				fmt.Println("Token:    expired or invalid")
			} else {
				fmt.Println("Token:    valid")
			}

			return nil
		},
	}
}

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internalauth.Logout(); err != nil {
				return exitcode.AuthError("could not logout: %v", err)
			}
			output.Success("Logged out successfully")
			return nil
		},
	}
}
