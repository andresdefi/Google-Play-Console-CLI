package apprecovery

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
		Use:     "apprecovery",
		Aliases: []string{"app-recovery", "recovery"},
		Short:   "Manage app recovery actions",
		Long:    "List, create, deploy, cancel, or add targeting to app recovery actions.",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newDeployCmd())
	cmd.AddCommand(newCancelCmd())
	cmd.AddCommand(newAddTargetingCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List app recovery actions",
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
			resp, err := client.Get(api.AppRecoveriesPath(pkg), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var list struct {
					RecoveryActions []struct {
						RecoveryID string `json:"appRecoveryId"`
						Status     string `json:"status"`
					} `json:"recoveryActions"`
				}
				raw := data.(json.RawMessage)
				if err := json.Unmarshal(raw, &list); err == nil && len(list.RecoveryActions) > 0 {
					t := output.NewTable(w, "Recovery ID", "Status")
					for _, r := range list.RecoveryActions {
						t.AppendRow([]any{r.RecoveryID, r.Status})
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

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create an app recovery action",
		Long:  "Create a new app recovery action. Pass the recovery data as JSON via stdin.",
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
				return exitcode.ConfigError("could not read recovery data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.AppRecoveriesPath(pkg), body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
}

func newDeployCmd() *cobra.Command {
	var recoveryID string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an app recovery action",
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
			_, err = client.Post(api.AppRecoveryPath(pkg, recoveryID)+":deploy", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("App recovery %s deployed", recoveryID))
			return nil
		},
	}
	cmd.Flags().StringVar(&recoveryID, "recovery-id", "", "Recovery action ID (required)")
	_ = cmd.MarkFlagRequired("recovery-id")
	return cmd
}

func newCancelCmd() *cobra.Command {
	var recoveryID string

	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel an app recovery action",
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
			_, err = client.Post(api.AppRecoveryPath(pkg, recoveryID)+":cancel", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("App recovery %s cancelled", recoveryID))
			return nil
		},
	}
	cmd.Flags().StringVar(&recoveryID, "recovery-id", "", "Recovery action ID (required)")
	_ = cmd.MarkFlagRequired("recovery-id")
	return cmd
}

func newAddTargetingCmd() *cobra.Command {
	var recoveryID string

	cmd := &cobra.Command{
		Use:   "add-targeting",
		Short: "Add targeting to an app recovery action",
		Long:  "Add targeting to an app recovery action. Pass the targeting data as JSON via stdin.",
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
				return exitcode.ConfigError("could not read targeting data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.AppRecoveryPath(pkg, recoveryID)+":addTargeting", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&recoveryID, "recovery-id", "", "Recovery action ID (required)")
	_ = cmd.MarkFlagRequired("recovery-id")
	return cmd
}
