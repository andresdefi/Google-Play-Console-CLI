package purchases

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
		Use:     "purchases",
		Aliases: []string{"purchase"},
		Short:   "Manage purchases",
		Long:    "Manage in-app product purchases, subscriptions, and voided purchases.",
	}

	cmd.AddCommand(newProductsCmd())
	cmd.AddCommand(newProductsV2Cmd())
	cmd.AddCommand(newSubscriptionsCmd())
	cmd.AddCommand(newSubscriptionsV2Cmd())
	cmd.AddCommand(newVoidedCmd())
	return cmd
}

// --- products ---

func newProductsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "products",
		Short: "Manage product purchases",
	}

	cmd.AddCommand(newProductsGetCmd())
	cmd.AddCommand(newProductsAcknowledgeCmd())
	cmd.AddCommand(newProductsConsumeCmd())
	return cmd
}

func newProductsGetCmd() *cobra.Command {
	var (
		productID string
		tok       string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a product purchase",
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
			resp, err := client.Get(api.PurchaseProductPath(pkg, productID, tok), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Product ID (required)")
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func newProductsAcknowledgeCmd() *cobra.Command {
	var (
		productID string
		tok       string
	)

	cmd := &cobra.Command{
		Use:   "acknowledge",
		Short: "Acknowledge a product purchase",
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
			_, err = client.Post(api.PurchaseProductPath(pkg, productID, tok)+":acknowledge", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Product purchase acknowledged for %s", productID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Product ID (required)")
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func newProductsConsumeCmd() *cobra.Command {
	var (
		productID string
		tok       string
	)

	cmd := &cobra.Command{
		Use:   "consume",
		Short: "Consume a product purchase",
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
			_, err = client.Post(api.PurchaseProductPath(pkg, productID, tok)+":consume", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Product purchase consumed for %s", productID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Product ID (required)")
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

// --- products-v2 ---

func newProductsV2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "products-v2",
		Short: "Manage product purchases (v2)",
	}

	cmd.AddCommand(newProductsV2GetCmd())
	return cmd
}

func newProductsV2GetCmd() *cobra.Command {
	var tok string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a product purchase (v2)",
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
			resp, err := client.Get(api.PurchaseProductV2Path(pkg, tok), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

// --- subscriptions ---

func newSubscriptionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "Manage subscription purchases",
	}

	cmd.AddCommand(newSubscriptionsGetCmd())
	cmd.AddCommand(newSubscriptionsAcknowledgeCmd())
	cmd.AddCommand(newSubscriptionsCancelCmd())
	cmd.AddCommand(newSubscriptionsDeferCmd())
	return cmd
}

func newSubscriptionsGetCmd() *cobra.Command {
	var (
		subscriptionID string
		tok            string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a subscription purchase",
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
			resp, err := client.Get(api.PurchaseSubscriptionPath(pkg, subscriptionID, tok), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "Subscription ID (required)")
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("subscription-id")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func newSubscriptionsAcknowledgeCmd() *cobra.Command {
	var (
		subscriptionID string
		tok            string
	)

	cmd := &cobra.Command{
		Use:   "acknowledge",
		Short: "Acknowledge a subscription purchase",
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
			_, err = client.Post(api.PurchaseSubscriptionPath(pkg, subscriptionID, tok)+":acknowledge", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Subscription purchase acknowledged for %s", subscriptionID))
			return nil
		},
	}
	cmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "Subscription ID (required)")
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("subscription-id")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func newSubscriptionsCancelCmd() *cobra.Command {
	var (
		subscriptionID string
		tok            string
	)

	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a subscription purchase",
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
			_, err = client.Post(api.PurchaseSubscriptionPath(pkg, subscriptionID, tok)+":cancel", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Subscription purchase cancelled for %s", subscriptionID))
			return nil
		},
	}
	cmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "Subscription ID (required)")
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("subscription-id")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func newSubscriptionsDeferCmd() *cobra.Command {
	var (
		subscriptionID string
		tok            string
	)

	cmd := &cobra.Command{
		Use:   "defer",
		Short: "Defer a subscription purchase",
		Long:  "Defer a subscription purchase. Pass the deferral data as JSON via stdin.",
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
				return exitcode.ConfigError("could not read deferral data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.PurchaseSubscriptionPath(pkg, subscriptionID, tok)+":defer", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&subscriptionID, "subscription-id", "", "Subscription ID (required)")
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("subscription-id")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

// --- subscriptions-v2 ---

func newSubscriptionsV2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions-v2",
		Short: "Manage subscription purchases (v2)",
	}

	cmd.AddCommand(newSubscriptionsV2GetCmd())
	cmd.AddCommand(newSubscriptionsV2CancelCmd())
	cmd.AddCommand(newSubscriptionsV2DeferCmd())
	cmd.AddCommand(newSubscriptionsV2RevokeCmd())
	return cmd
}

func newSubscriptionsV2GetCmd() *cobra.Command {
	var tok string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a subscription purchase (v2)",
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
			resp, err := client.Get(api.PurchaseSubscriptionV2Path(pkg, tok), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func newSubscriptionsV2CancelCmd() *cobra.Command {
	var tok string

	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a subscription purchase (v2)",
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
			_, err = client.Post(api.PurchaseSubscriptionV2Path(pkg, tok)+":cancel", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success("Subscription purchase cancelled (v2)")
			return nil
		},
	}
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func newSubscriptionsV2DeferCmd() *cobra.Command {
	var tok string

	cmd := &cobra.Command{
		Use:   "defer",
		Short: "Defer a subscription purchase (v2)",
		Long:  "Defer a subscription purchase (v2). Pass the deferral data as JSON via stdin.",
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
				return exitcode.ConfigError("could not read deferral data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.PurchaseSubscriptionV2Path(pkg, tok)+":defer", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func newSubscriptionsV2RevokeCmd() *cobra.Command {
	var tok string

	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a subscription purchase (v2)",
		Long:  "Revoke a subscription purchase (v2). Pass the revocation data as JSON via stdin.",
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
				return exitcode.ConfigError("could not read revocation data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.PurchaseSubscriptionV2Path(pkg, tok)+":revoke", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&tok, "token", "", "Purchase token (required)")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

// --- voided ---

func newVoidedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "voided",
		Short: "Manage voided purchases",
	}

	cmd.AddCommand(newVoidedListCmd())
	return cmd
}

func newVoidedListCmd() *cobra.Command {
	var (
		startTime string
		endTime   string
		token     string
		maxResult string
		pType     string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List voided purchases",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			authToken, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			params := map[string]string{}
			if startTime != "" {
				params["startTime"] = startTime
			}
			if endTime != "" {
				params["endTime"] = endTime
			}
			if token != "" {
				params["token"] = token
			}
			if maxResult != "" {
				params["maxResults"] = maxResult
			}
			if pType != "" {
				params["type"] = pType
			}

			client := api.NewClient(authToken)
			resp, err := client.Get(api.VoidedPurchasesPath(pkg), params)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var list struct {
					VoidedPurchases []struct {
						PurchaseToken    string `json:"purchaseToken"`
						PurchaseTimeMS   string `json:"purchaseTimeMillis"`
						VoidedTimeMS     string `json:"voidedTimeMillis"`
						OrderID          string `json:"orderId"`
						VoidedSource     int    `json:"voidedSource"`
						VoidedReason     int    `json:"voidedReason"`
					} `json:"voidedPurchases"`
				}
				raw := data.(json.RawMessage)
				if err := json.Unmarshal(raw, &list); err == nil && len(list.VoidedPurchases) > 0 {
					t := output.NewTable(w, "Order ID", "Purchase Time", "Voided Time", "Source", "Reason")
					for _, vp := range list.VoidedPurchases {
						t.AppendRow([]any{vp.OrderID, vp.PurchaseTimeMS, vp.VoidedTimeMS, vp.VoidedSource, vp.VoidedReason})
					}
					t.Render()
				} else {
					_, _ = fmt.Fprintln(w, string(raw))
				}
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&startTime, "start-time", "", "Start time in milliseconds since epoch")
	cmd.Flags().StringVar(&endTime, "end-time", "", "End time in milliseconds since epoch")
	cmd.Flags().StringVar(&token, "token", "", "Pagination token")
	cmd.Flags().StringVar(&maxResult, "max-results", "", "Max results per page")
	cmd.Flags().StringVar(&pType, "type", "", "Purchase type filter")
	return cmd
}
