package apps

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
		Use:     "apps",
		Aliases: []string{"app"},
		Short:   "Manage apps",
	}

	cmd.AddCommand(newGetCmd())
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
			// Use edits to get app details.
			edit, err := client.CreateEdit(pkg)
			if err != nil {
				return exitcode.APIErrorExit("could not create edit: %v", err)
			}
			defer client.DeleteEdit(pkg, edit.ID)

			resp, err := client.Get(api.DetailsPath(pkg, edit.ID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
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
