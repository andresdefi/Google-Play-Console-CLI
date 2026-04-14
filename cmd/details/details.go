package details

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
		Use:   "details",
		Short: "Manage app details",
		Long:  "Get or update app details (contact info, default language).",
	}

	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newUpdateCmd())
	return cmd
}

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get app details",
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
			resp, err := withTempEdit(client, pkg, func(editID string) (json.RawMessage, error) {
				return client.Get(api.DetailsPath(pkg, editID), nil)
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, func(w io.Writer, data any) {
				var d struct {
					DefaultLanguage string `json:"defaultLanguage"`
					ContactEmail    string `json:"contactEmail"`
					ContactPhone    string `json:"contactPhone"`
					ContactWebsite  string `json:"contactWebsite"`
				}
				if err := json.Unmarshal(data.(json.RawMessage), &d); err == nil {
					t := output.NewTable(w, "Field", "Value")
					t.AppendRow([]any{"Package", pkg})
					t.AppendRow([]any{"Default Language", d.DefaultLanguage})
					t.AppendRow([]any{"Contact Email", d.ContactEmail})
					t.AppendRow([]any{"Contact Phone", d.ContactPhone})
					t.AppendRow([]any{"Contact Website", d.ContactWebsite})
					t.Render()
				} else {
					fmt.Fprintln(w, string(data.(json.RawMessage)))
				}
			})
			return nil
		},
	}
}

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update app details",
		Long:  "Update app details. Pass the details data as JSON via stdin.",
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
				var body json.RawMessage
				if err := json.NewDecoder(cmd.InOrStdin()).Decode(&body); err != nil {
					return fmt.Errorf("could not read details data from stdin: %w", err)
				}
				_, err := client.Put(api.DetailsPath(pkg, editID), body)
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success("App details updated and committed")
			return nil
		},
	}
}

// withTempEdit creates a temporary read-only edit, runs the function, then deletes the edit.
func withTempEdit(client *api.Client, pkg string, fn func(editID string) (json.RawMessage, error)) (json.RawMessage, error) {
	edit, err := client.CreateEdit(pkg)
	if err != nil {
		return nil, err
	}
	defer client.DeleteEdit(pkg, edit.ID)

	return fn(edit.ID)
}
