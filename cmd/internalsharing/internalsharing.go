package internalsharing

import (
	"encoding/json"
	"os"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "internalsharing",
		Aliases: []string{"internal-sharing", "share"},
		Short:   "Manage internal app sharing",
		Long:    "Upload APKs or bundles for internal app sharing.",
	}

	cmd.AddCommand(newUploadAPKCmd())
	cmd.AddCommand(newUploadBundleCmd())
	return cmd
}

func newUploadAPKCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload-apk <file>",
		Short: "Upload an APK for internal sharing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			if _, err := os.Stat(filePath); err != nil {
				return exitcode.NewExitError(exitcode.Error, "file not found: %s", filePath)
			}

			client := api.NewClient(token)
			resp, err := client.Upload(api.InternalSharingAPKPath(pkg), filePath, "application/vnd.android.package-archive")
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
}

func newUploadBundleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload-bundle <file>",
		Short: "Upload a bundle for internal sharing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			if _, err := os.Stat(filePath); err != nil {
				return exitcode.NewExitError(exitcode.Error, "file not found: %s", filePath)
			}

			client := api.NewClient(token)
			resp, err := client.Upload(api.InternalSharingBundlePath(pkg), filePath, "application/octet-stream")
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
}
