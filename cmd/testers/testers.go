package testers

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
		Use:     "testers",
		Aliases: []string{"tester"},
		Short:   "Manage testers",
		Long:    "Get or update testers for a release track.",
	}

	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newUpdateCmd())
	return cmd
}

func newGetCmd() *cobra.Command {
	var track string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get testers for a track",
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
				return client.Get(api.TestersPath(pkg, editID, track), nil)
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, func(w io.Writer, data any) {
				var result struct {
					GoogleGroups []string `json:"googleGroups"`
				}
				if err := json.Unmarshal(data.(json.RawMessage), &result); err == nil && len(result.GoogleGroups) > 0 {
					t := output.NewTable(w, "Google Group")
					for _, g := range result.GoogleGroups {
						t.AppendRow([]any{g})
					}
					t.Render()
				} else {
					fmt.Fprintln(w, string(data.(json.RawMessage)))
				}
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&track, "track", "", "Track name (required)")
	cmd.MarkFlagRequired("track")
	return cmd
}

func newUpdateCmd() *cobra.Command {
	var track string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update testers for a track",
		Long:  "Update testers for a track. Pass the testers data as JSON via stdin.",
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
					return fmt.Errorf("could not read testers data from stdin: %w", err)
				}
				_, err := client.Put(api.TestersPath(pkg, editID, track), body)
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Testers for %s updated and committed", track))
			return nil
		},
	}
	cmd.Flags().StringVar(&track, "track", "", "Track name (required)")
	cmd.MarkFlagRequired("track")
	return cmd
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
