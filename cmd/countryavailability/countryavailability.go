package countryavailability

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
		Use:     "countryavailability",
		Aliases: []string{"country-availability", "countries"},
		Short:   "Manage country availability",
		Long:    "Get country availability for a track.",
	}

	cmd.AddCommand(newGetCmd())
	return cmd
}

func newGetCmd() *cobra.Command {
	var track string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get country availability for a track",
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

			resp, err := client.Get(api.CountryAvailabilityPath(pkg, edit.ID, track), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var ca struct {
					Countries []struct {
						CountryCode string `json:"countryCode"`
					} `json:"countries"`
				}
				raw := data.(json.RawMessage)
				if err := json.Unmarshal(raw, &ca); err == nil && len(ca.Countries) > 0 {
					t := output.NewTable(w, "Country Code")
					for _, c := range ca.Countries {
						t.AppendRow([]any{c.CountryCode})
					}
					t.Render()
				} else {
					fmt.Fprintln(w, string(raw))
				}
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&track, "track", "", "Track name (required)")
	cmd.MarkFlagRequired("track")
	return cmd
}
