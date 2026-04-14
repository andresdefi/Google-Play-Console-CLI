package edits

import (
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
		Use:   "edits",
		Short: "Manage edit sessions",
		Long:  "Create, validate, commit, or delete edit sessions. Edits are transactional containers for app changes.",
	}

	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newCommitCmd())
	cmd.AddCommand(newDeleteCmd())
	return cmd
}

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a new edit session",
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

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, edit, func(w io.Writer, data any) {
				e := data.(*api.Edit)
				t := output.NewTable(w, "Edit ID", "Expiry")
				t.AppendRow([]any{e.ID, e.ExpiryTimeSeconds})
				t.Render()
			})
			return nil
		},
	}
}

func newGetCmd() *cobra.Command {
	var editID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an edit session",
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
			edit, err := client.GetEdit(pkg, editID)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, edit, func(w io.Writer, data any) {
				e := data.(*api.Edit)
				t := output.NewTable(w, "Edit ID", "Expiry")
				t.AppendRow([]any{e.ID, e.ExpiryTimeSeconds})
				t.Render()
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&editID, "edit-id", "", "Edit session ID (required)")
	_ = cmd.MarkFlagRequired("edit-id")
	return cmd
}

func newValidateCmd() *cobra.Command {
	var editID string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate an edit session",
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
			_, err = client.ValidateEdit(pkg, editID)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Edit %s is valid", editID))
			return nil
		},
	}
	cmd.Flags().StringVar(&editID, "edit-id", "", "Edit session ID (required)")
	_ = cmd.MarkFlagRequired("edit-id")
	return cmd
}

func newCommitCmd() *cobra.Command {
	var editID string

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Commit an edit session",
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
			committed, err := client.CommitEdit(pkg, editID)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			_ = committed
			output.Success(fmt.Sprintf("Edit %s committed successfully", editID))
			return nil
		},
	}
	cmd.Flags().StringVar(&editID, "edit-id", "", "Edit session ID (required)")
	_ = cmd.MarkFlagRequired("edit-id")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	var editID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an edit session",
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
			if err := client.DeleteEdit(pkg, editID); err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Edit %s deleted", editID))
			return nil
		},
	}
	cmd.Flags().StringVar(&editID, "edit-id", "", "Edit session ID (required)")
	_ = cmd.MarkFlagRequired("edit-id")
	return cmd
}

