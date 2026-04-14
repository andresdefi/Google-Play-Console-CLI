package orders

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
		Use:     "orders",
		Aliases: []string{"order"},
		Short:   "Manage orders",
		Long:    "Get, batch-get, or refund orders for an app.",
	}

	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newBatchGetCmd())
	cmd.AddCommand(newRefundCmd())
	return cmd
}

func newGetCmd() *cobra.Command {
	var orderID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an order",
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
			resp, err := client.Get(api.OrderPath(pkg, orderID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&orderID, "order-id", "", "Order ID (required)")
	_ = cmd.MarkFlagRequired("order-id")
	return cmd
}

func newBatchGetCmd() *cobra.Command {
	var orderIDs string

	cmd := &cobra.Command{
		Use:   "batch-get",
		Short: "Batch get orders",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			ids := strings.Split(orderIDs, ",")
			params := map[string]string{}
			for i, id := range ids {
				params[fmt.Sprintf("orderIds[%d]", i)] = strings.TrimSpace(id)
			}

			client := api.NewClient(token)
			resp, err := client.Get(api.OrdersPath(pkg)+":batchGet", params)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&orderIDs, "order-ids", "", "Comma-separated list of order IDs (required)")
	_ = cmd.MarkFlagRequired("order-ids")
	return cmd
}

func newRefundCmd() *cobra.Command {
	var orderID string

	cmd := &cobra.Command{
		Use:   "refund",
		Short: "Refund an order",
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
			_, err = client.Post(api.OrderPath(pkg, orderID)+":refund", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Order %s refunded", orderID))
			return nil
		},
	}
	cmd.Flags().StringVar(&orderID, "order-id", "", "Order ID (required)")
	_ = cmd.MarkFlagRequired("order-id")
	return cmd
}
