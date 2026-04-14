package generatedapks

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generatedapks",
		Aliases: []string{"generated-apk"},
		Short:   "Manage generated APKs",
		Long:    "List or download generated APKs for an app.",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newDownloadCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	var versionCode int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List generated APKs",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Get(api.GeneratedAPKsPath(pkg, versionCode), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var list struct {
					GeneratedAPKs []struct {
						DownloadID    string `json:"downloadId"`
						VariantID     int    `json:"variantId"`
						CertSHA256    string `json:"certificateSha256Fingerprint"`
					} `json:"generatedApks"`
				}
				raw := data.(json.RawMessage)
				if err := json.Unmarshal(raw, &list); err == nil && len(list.GeneratedAPKs) > 0 {
					t := output.NewTable(w, "Download ID", "Variant ID", "Cert SHA256")
					for _, a := range list.GeneratedAPKs {
						t.AppendRow([]any{a.DownloadID, a.VariantID, a.CertSHA256})
					}
					t.Render()
				} else {
					fmt.Fprintln(w, string(raw))
				}
			})
			return nil
		},
	}
	cmd.Flags().IntVar(&versionCode, "version-code", 0, "Version code (required)")
	cmd.MarkFlagRequired("version-code")
	return cmd
}

func newDownloadCmd() *cobra.Command {
	var (
		versionCode int
		downloadID  string
		outputPath  string
	)

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download a generated APK",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			client := api.NewClient(token)
			err = client.DownloadToFile(api.GeneratedAPKDownloadPath(pkg, versionCode, downloadID), outputPath)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Downloaded to %s", outputPath))
			return nil
		},
	}
	cmd.Flags().IntVar(&versionCode, "version-code", 0, "Version code (required)")
	cmd.Flags().StringVar(&downloadID, "download-id", "", "Download ID (required)")
	cmd.Flags().StringVar(&outputPath, "output", "", "Output file path (required)")
	cmd.MarkFlagRequired("version-code")
	cmd.MarkFlagRequired("download-id")
	cmd.MarkFlagRequired("output")
	return cmd
}
