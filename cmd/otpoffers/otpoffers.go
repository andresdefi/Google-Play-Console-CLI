package otpoffers

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
		Use:     "otpoffers",
		Aliases: []string{"otp-offer"},
		Short:   "Manage one-time product offers",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newActivateCmd())
	cmd.AddCommand(newDeactivateCmd())
	cmd.AddCommand(newCancelCmd())
	cmd.AddCommand(newBatchGetCmd())
	cmd.AddCommand(newBatchUpdateCmd())
	cmd.AddCommand(newBatchDeleteCmd())
	cmd.AddCommand(newBatchUpdateStatesCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	var (
		productID        string
		purchaseOptionID string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List one-time product offers",
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
			resp, err := client.Get(api.OTPOffersPath(pkg, productID, purchaseOptionID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.Flags().StringVar(&purchaseOptionID, "purchase-option-id", "", "Purchase option ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("purchase-option-id")
	return cmd
}

func newActivateCmd() *cobra.Command {
	var (
		productID        string
		purchaseOptionID string
		offerID          string
	)

	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Activate a one-time product offer",
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
			_, err = client.Post(api.OTPOfferPath(pkg, productID, purchaseOptionID, offerID)+":activate", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("OTP offer %s activated", offerID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.Flags().StringVar(&purchaseOptionID, "purchase-option-id", "", "Purchase option ID (required)")
	cmd.Flags().StringVar(&offerID, "offer-id", "", "Offer ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("purchase-option-id")
	_ = cmd.MarkFlagRequired("offer-id")
	return cmd
}

func newDeactivateCmd() *cobra.Command {
	var (
		productID        string
		purchaseOptionID string
		offerID          string
	)

	cmd := &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate a one-time product offer",
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
			_, err = client.Post(api.OTPOfferPath(pkg, productID, purchaseOptionID, offerID)+":deactivate", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("OTP offer %s deactivated", offerID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.Flags().StringVar(&purchaseOptionID, "purchase-option-id", "", "Purchase option ID (required)")
	cmd.Flags().StringVar(&offerID, "offer-id", "", "Offer ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("purchase-option-id")
	_ = cmd.MarkFlagRequired("offer-id")
	return cmd
}

func newCancelCmd() *cobra.Command {
	var (
		productID        string
		purchaseOptionID string
		offerID          string
	)

	cmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a one-time product offer",
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
			_, err = client.Post(api.OTPOfferPath(pkg, productID, purchaseOptionID, offerID)+":cancel", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("OTP offer %s cancelled", offerID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.Flags().StringVar(&purchaseOptionID, "purchase-option-id", "", "Purchase option ID (required)")
	cmd.Flags().StringVar(&offerID, "offer-id", "", "Offer ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("purchase-option-id")
	_ = cmd.MarkFlagRequired("offer-id")
	return cmd
}

func newBatchGetCmd() *cobra.Command {
	var (
		productID        string
		purchaseOptionID string
	)

	cmd := &cobra.Command{
		Use:   "batch-get",
		Short: "Batch get one-time product offers",
		Long:  "Batch get one-time product offers. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.OTPOffersPath(pkg, productID, purchaseOptionID)+":batchGet", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.Flags().StringVar(&purchaseOptionID, "purchase-option-id", "", "Purchase option ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("purchase-option-id")
	return cmd
}

func newBatchUpdateCmd() *cobra.Command {
	var (
		productID        string
		purchaseOptionID string
	)

	cmd := &cobra.Command{
		Use:   "batch-update",
		Short: "Batch update one-time product offers",
		Long:  "Batch update one-time product offers. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.OTPOffersPath(pkg, productID, purchaseOptionID)+":batchUpdate", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.Flags().StringVar(&purchaseOptionID, "purchase-option-id", "", "Purchase option ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("purchase-option-id")
	return cmd
}

func newBatchDeleteCmd() *cobra.Command {
	var (
		productID        string
		purchaseOptionID string
	)

	cmd := &cobra.Command{
		Use:   "batch-delete",
		Short: "Batch delete one-time product offers",
		Long:  "Batch delete one-time product offers. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.OTPOffersPath(pkg, productID, purchaseOptionID)+":batchDelete", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.Flags().StringVar(&purchaseOptionID, "purchase-option-id", "", "Purchase option ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("purchase-option-id")
	return cmd
}

func newBatchUpdateStatesCmd() *cobra.Command {
	var (
		productID        string
		purchaseOptionID string
	)

	cmd := &cobra.Command{
		Use:   "batch-update-states",
		Short: "Batch update one-time product offer states",
		Long:  "Batch update one-time product offer states. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.OTPOffersPath(pkg, productID, purchaseOptionID)+":batchUpdateStates", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.Flags().StringVar(&purchaseOptionID, "purchase-option-id", "", "Purchase option ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("purchase-option-id")
	return cmd
}
