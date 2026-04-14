package deobfuscation

import (
	"fmt"
	"os"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deobfuscation",
		Short: "Manage deobfuscation files",
		Long:  "Upload ProGuard mapping or native debug symbol files for an APK.",
	}

	cmd.AddCommand(newUploadCmd())
	return cmd
}

func newUploadCmd() *cobra.Command {
	var (
		versionCode int
		fileType    string
	)

	cmd := &cobra.Command{
		Use:   "upload <file>",
		Short: "Upload a deobfuscation file",
		Long:  "Upload a ProGuard mapping file or native debug symbols for a specific APK version.",
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
			_, err = client.WithEdit(pkg, func(editID string) error {
				_, err := client.Upload(api.DeobfuscationFilesPath(pkg, editID, versionCode, fileType), filePath, "application/octet-stream")
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Deobfuscation file uploaded and committed for %s (version %d, type %s)", pkg, versionCode, fileType))
			return nil
		},
	}
	cmd.Flags().IntVar(&versionCode, "apk-version", 0, "APK version code (required)")
	cmd.Flags().StringVar(&fileType, "type", "", "File type: proguard or nativeCode (required)")
	_ = cmd.MarkFlagRequired("apk-version")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}
