package users

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "users",
		Aliases: []string{"user"},
		Short:   "Manage users",
		Long:    "List, create, update, or delete developer account users.",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newDeleteCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	var developerID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Get(api.UsersPath(developerID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var list struct {
					Users []struct {
						Name  string `json:"name"`
						Email string `json:"email"`
					} `json:"users"`
				}
				raw := data.(json.RawMessage)
				if err := json.Unmarshal(raw, &list); err == nil && len(list.Users) > 0 {
					t := output.NewTable(w, "Name", "Email")
					for _, u := range list.Users {
						t.AppendRow([]any{u.Name, u.Email})
					}
					t.Render()
				} else {
					_, _ = fmt.Fprintln(w, string(raw))
				}
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&developerID, "developer-id", "", "Developer account ID (required)")
	_ = cmd.MarkFlagRequired("developer-id")
	return cmd
}

func newCreateCmd() *cobra.Command {
	var developerID string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a user",
		Long:  "Create a new user. Pass the user data as JSON via stdin.",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			var body json.RawMessage
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&body); err != nil {
				return exitcode.ConfigError("could not read user data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.UsersPath(developerID), body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&developerID, "developer-id", "", "Developer account ID (required)")
	_ = cmd.MarkFlagRequired("developer-id")
	return cmd
}

func newUpdateCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a user",
		Long:  "Update a user. Pass the user data as JSON via stdin. The --name flag is the resource name (developers/{devId}/users/{userId}).",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			var body json.RawMessage
			if err := json.NewDecoder(cmd.InOrStdin()).Decode(&body); err != nil {
				return exitcode.ConfigError("could not read user data from stdin: %v", err)
			}

			// Parse developer ID and user ID from resource name.
			parts := strings.Split(name, "/")
			if len(parts) < 4 || parts[0] != "developers" || parts[2] != "users" {
				return exitcode.ConfigError("invalid resource name: %s (expected developers/{devId}/users/{userId})", name)
			}
			developerID := parts[1]
			userID := parts[3]

			client := api.NewClient(token)
			resp, err := client.Patch(api.UserPath(developerID, userID), nil, body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Resource name (developers/{devId}/users/{userId}) (required)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			// Parse developer ID and user ID from resource name.
			parts := strings.Split(name, "/")
			if len(parts) < 4 || parts[0] != "developers" || parts[2] != "users" {
				return exitcode.ConfigError("invalid resource name: %s (expected developers/{devId}/users/{userId})", name)
			}
			developerID := parts[1]
			userID := parts[3]

			client := api.NewClient(token)
			err = client.Delete(api.UserPath(developerID, userID))
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("User %s deleted", name))
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Resource name (developers/{devId}/users/{userId}) (required)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
