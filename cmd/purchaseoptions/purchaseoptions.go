package purchaseoptions

import (
	"encoding/json"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "purchaseoptions",
		Aliases: []string{"purchase-option", "po"},
		Short:   "Manage purchase options",
	}

	cmd.AddCommand(newBatchDeleteCmd())
	cmd.AddCommand(newBatchUpdateStatesCmd())
	return cmd
}

func newBatchDeleteCmd() *cobra.Command {
	var productID string

	cmd := &cobra.Command{
		Use:   "batch-delete",
		Short: "Batch delete purchase options",
		Long:  "Batch delete purchase options. Pass the request body as JSON via stdin.",
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
				return exitcode.ConfigError("could not read JSON from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.PurchaseOptionsPath(pkg, productID)+":batchDelete", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.MarkFlagRequired("product-id")
	return cmd
}

func newBatchUpdateStatesCmd() *cobra.Command {
	var productID string

	cmd := &cobra.Command{
		Use:   "batch-update-states",
		Short: "Batch update purchase option states",
		Long:  "Batch update purchase option states. Pass the request body as JSON via stdin.",
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
				return exitcode.ConfigError("could not read JSON from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.PurchaseOptionsPath(pkg, productID)+":batchUpdateStates", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.MarkFlagRequired("product-id")
	return cmd
}
