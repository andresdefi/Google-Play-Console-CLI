package grants

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "grants",
		Aliases: []string{"grant"},
		Short:   "Manage grants",
		Long:    "Create, update, or delete user grants.",
	}

	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newDeleteCmd())
	return cmd
}

func newCreateCmd() *cobra.Command {
	var (
		developerID string
		userID      string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a grant",
		Long:  "Create a new grant. Pass the grant data as JSON via stdin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			var body json.RawMessage
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&body); err != nil {
				return exitcode.ConfigError("could not read grant data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.GrantsPath(developerID, userID), body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&developerID, "developer-id", "", "Developer account ID (required)")
	cmd.Flags().StringVar(&userID, "user-id", "", "User ID (required)")
	cmd.MarkFlagRequired("developer-id")
	cmd.MarkFlagRequired("user-id")
	return cmd
}

func newUpdateCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a grant",
		Long:  "Update a grant. Pass the grant data as JSON via stdin. The --name flag is the resource name (developers/{devId}/users/{userId}/grants/{grantId}).",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			var body json.RawMessage
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&body); err != nil {
				return exitcode.ConfigError("could not read grant data from stdin: %v", err)
			}

			// Parse developer ID, user ID, and grant ID from resource name.
			parts := strings.Split(name, "/")
			if len(parts) < 6 || parts[0] != "developers" || parts[2] != "users" || parts[4] != "grants" {
				return exitcode.ConfigError("invalid resource name: %s (expected developers/{devId}/users/{userId}/grants/{grantId})", name)
			}
			developerID := parts[1]
			userID := parts[3]
			grantID := parts[5]

			client := api.NewClient(token)
			resp, err := client.Patch(api.GrantPath(developerID, userID, grantID), nil, body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Resource name (developers/{devId}/users/{userId}/grants/{grantId}) (required)")
	cmd.MarkFlagRequired("name")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			// Parse developer ID, user ID, and grant ID from resource name.
			parts := strings.Split(name, "/")
			if len(parts) < 6 || parts[0] != "developers" || parts[2] != "users" || parts[4] != "grants" {
				return exitcode.ConfigError("invalid resource name: %s (expected developers/{devId}/users/{userId}/grants/{grantId})", name)
			}
			developerID := parts[1]
			userID := parts[3]
			grantID := parts[5]

			client := api.NewClient(token)
			err = client.Delete(api.GrantPath(developerID, userID, grantID))
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Grant %s deleted", name))
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Resource name (developers/{devId}/users/{userId}/grants/{grantId}) (required)")
	cmd.MarkFlagRequired("name")
	return cmd
}
