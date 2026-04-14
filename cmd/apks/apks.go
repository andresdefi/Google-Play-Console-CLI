package apks

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apks",
		Aliases: []string{"apk"},
		Short:   "Manage APKs",
		Long:    "List, upload, or add externally hosted APKs for an app.",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newUploadCmd())
	cmd.AddCommand(newAddExternallyHostedCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all APKs",
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
			edit, err := client.CreateEdit(pkg)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}
			defer client.DeleteEdit(pkg, edit.ID)

			resp, err := client.Get(api.APKsPath(pkg, edit.ID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var list struct {
					Kind string `json:"kind"`
					APKs []struct {
						VersionCode int    `json:"versionCode"`
						Binary      struct {
							SHA1    string `json:"sha1"`
							SHA256  string `json:"sha256"`
						} `json:"binary"`
					} `json:"apks"`
				}
				raw := data.(json.RawMessage)
				if err := json.Unmarshal(raw, &list); err == nil && len(list.APKs) > 0 {
					t := output.NewTable(w, "Version Code", "SHA1", "SHA256")
					for _, a := range list.APKs {
						t.AppendRow([]any{a.VersionCode, a.Binary.SHA1, a.Binary.SHA256})
					}
					t.Render()
				} else {
					fmt.Fprintln(w, string(raw))
				}
			})
			return nil
		},
	}
}

func newUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload <file>",
		Short: "Upload an APK",
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
				_, err := client.Upload(api.APKsPath(pkg, editID), filePath, "application/vnd.android.package-archive")
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("APK uploaded and committed for %s", pkg))
			return nil
		},
	}
}

func newAddExternallyHostedCmd() *cobra.Command {
	var url string

	cmd := &cobra.Command{
		Use:   "add-externally-hosted",
		Short: "Add an externally hosted APK",
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
			_, err = client.WithEdit(pkg, func(editID string) error {
				body := map[string]any{
					"externallyHostedUrl": url,
				}
				_, err := client.Post(api.APKsPath(pkg, editID)+"/externallyHosted", body)
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Externally hosted APK added and committed for %s", pkg))
			return nil
		},
	}
	cmd.Flags().StringVar(&url, "url", "", "URL of the externally hosted APK (required)")
	cmd.MarkFlagRequired("url")
	return cmd
}
