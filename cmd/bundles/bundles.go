package bundles

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bundles",
		Aliases: []string{"bundle"},
		Short:   "Manage app bundles",
		Long:    "List or upload Android App Bundles (AAB) for an app.",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newUploadCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all bundles",
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
			defer func() { _ = client.DeleteEdit(pkg, edit.ID) }()

			resp, err := client.Get(api.BundlesPath(pkg, edit.ID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var list struct {
					Kind    string `json:"kind"`
					Bundles []struct {
						VersionCode int    `json:"versionCode"`
						SHA1        string `json:"sha1"`
						SHA256      string `json:"sha256"`
					} `json:"bundles"`
				}
				raw := data.(json.RawMessage)
				if err := json.Unmarshal(raw, &list); err == nil && len(list.Bundles) > 0 {
					t := output.NewTable(w, "Version Code", "SHA1", "SHA256")
					for _, b := range list.Bundles {
						t.AppendRow([]any{b.VersionCode, b.SHA1, b.SHA256})
					}
					t.Render()
				} else {
					_, _ = fmt.Fprintln(w, string(raw))
				}
			})
			return nil
		},
	}
}

func newUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload <file>",
		Short: "Upload an app bundle",
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
				_, err := client.Upload(api.BundlesPath(pkg, editID), filePath, "application/octet-stream")
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", withUploadGuidance(err))
			}

			output.Success(fmt.Sprintf("Bundle uploaded and committed for %s", pkg))
			return nil
		},
	}
}

func withUploadGuidance(err error) error {
	var apiErr *api.APIError
	if !errors.As(err, &apiErr) {
		return err
	}
	if !shouldSuggestPlayConsoleSetup(apiErr) {
		return err
	}
	return fmt.Errorf("%w\n\nUpload failed - new apps may require completing setup in the Play Console first (store listing, content rating, data safety form). Also verify the service account has been invited in Play Console > Users and permissions for this app. Visit https://play.google.com/console to complete setup, then retry", err)
}

func shouldSuggestPlayConsoleSetup(err *api.APIError) bool {
	if err == nil {
		return false
	}
	if err.StatusCode == 404 || err.StatusCode == 403 {
		return true
	}
	text := strings.ToLower(err.Error())
	return strings.Contains(text, "not found") ||
		strings.Contains(text, "permission") ||
		strings.Contains(text, "access") ||
		strings.Contains(text, "not configured") ||
		strings.Contains(text, "setup")
}
