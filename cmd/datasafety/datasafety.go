package datasafety

import (
	"encoding/json"
	"fmt"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datasafety",
		Aliases: []string{"data-safety"},
		Short:   "Manage data safety declarations",
		Long:    "Update data safety declarations for an app.",
	}

	cmd.AddCommand(newUpdateCmd())
	return cmd
}

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update data safety declarations",
		Long:  "Update data safety declarations. Pass the data safety data as JSON via stdin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			var body json.RawMessage
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&body); err != nil {
				return exitcode.APIErrorExit("could not read data safety data from stdin: %v", err)
			}

			client := api.NewClient(token)
			_, err = client.Post(api.DataSafetyPath(pkg), body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Data safety declarations updated for %s", pkg))
			return nil
		},
	}
}
