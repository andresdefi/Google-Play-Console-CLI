package externaltransactions

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
		Use:     "externaltransactions",
		Aliases: []string{"external-transaction", "ext-txn"},
		Short:   "Manage external transactions",
		Long:    "Create, get, or refund external transactions for an app.",
	}

	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newRefundCmd())
	return cmd
}

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create an external transaction",
		Long:  "Create a new external transaction. Pass the transaction data as JSON via stdin.",
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
				return exitcode.ConfigError("could not read transaction data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.ExternalTransactionsPath(pkg), body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
}

func newGetCmd() *cobra.Command {
	var txnID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an external transaction",
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
			resp, err := client.Get(api.ExternalTransactionPath(pkg, txnID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&txnID, "transaction-id", "", "External transaction ID (required)")
	cmd.MarkFlagRequired("transaction-id")
	return cmd
}

func newRefundCmd() *cobra.Command {
	var txnID string

	cmd := &cobra.Command{
		Use:   "refund",
		Short: "Refund an external transaction",
		Long:  "Refund an external transaction. Pass the refund data as JSON via stdin.",
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
				return exitcode.ConfigError("could not read refund data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.ExternalTransactionPath(pkg, txnID)+":refund", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&txnID, "transaction-id", "", "External transaction ID (required)")
	cmd.MarkFlagRequired("transaction-id")
	return cmd
}
