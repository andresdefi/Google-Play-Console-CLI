package offers

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
		Use:     "offers",
		Aliases: []string{"offer"},
		Short:   "Manage subscription offers",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newActivateCmd())
	cmd.AddCommand(newDeactivateCmd())
	cmd.AddCommand(newBatchGetCmd())
	cmd.AddCommand(newBatchUpdateCmd())
	cmd.AddCommand(newBatchUpdateStatesCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List offers",
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
			resp, err := client.Get(api.OffersPath(pkg, productID, basePlanID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("base-plan-id")
	return cmd
}

func newGetCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
		offerID    string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an offer",
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
			resp, err := client.Get(api.OfferPath(pkg, productID, basePlanID, offerID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	cmd.Flags().StringVar(&offerID, "offer-id", "", "Offer ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("base-plan-id")
	_ = cmd.MarkFlagRequired("offer-id")
	return cmd
}

func newCreateCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an offer",
		Long:  "Create an offer. Pass the offer data as JSON via stdin.",
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
			resp, err := client.Post(api.OffersPath(pkg, productID, basePlanID), body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("base-plan-id")
	return cmd
}

func newUpdateCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
		offerID    string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an offer",
		Long:  "Update an offer. Pass the offer data as JSON via stdin.",
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
			resp, err := client.Patch(api.OfferPath(pkg, productID, basePlanID, offerID), nil, body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	cmd.Flags().StringVar(&offerID, "offer-id", "", "Offer ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("base-plan-id")
	_ = cmd.MarkFlagRequired("offer-id")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
		offerID    string
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an offer",
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
			if err := client.Delete(api.OfferPath(pkg, productID, basePlanID, offerID)); err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Offer %s deleted", offerID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	cmd.Flags().StringVar(&offerID, "offer-id", "", "Offer ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("base-plan-id")
	_ = cmd.MarkFlagRequired("offer-id")
	return cmd
}

func newActivateCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
		offerID    string
	)

	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Activate an offer",
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
			_, err = client.Post(api.OfferPath(pkg, productID, basePlanID, offerID)+":activate", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Offer %s activated", offerID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	cmd.Flags().StringVar(&offerID, "offer-id", "", "Offer ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("base-plan-id")
	_ = cmd.MarkFlagRequired("offer-id")
	return cmd
}

func newDeactivateCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
		offerID    string
	)

	cmd := &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate an offer",
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
			_, err = client.Post(api.OfferPath(pkg, productID, basePlanID, offerID)+":deactivate", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Offer %s deactivated", offerID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	cmd.Flags().StringVar(&offerID, "offer-id", "", "Offer ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("base-plan-id")
	_ = cmd.MarkFlagRequired("offer-id")
	return cmd
}

func newBatchGetCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
	)

	cmd := &cobra.Command{
		Use:   "batch-get",
		Short: "Batch get offers",
		Long:  "Batch get offers. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.OffersPath(pkg, productID, basePlanID)+":batchGet", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("base-plan-id")
	return cmd
}

func newBatchUpdateCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
	)

	cmd := &cobra.Command{
		Use:   "batch-update",
		Short: "Batch update offers",
		Long:  "Batch update offers. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.OffersPath(pkg, productID, basePlanID)+":batchUpdate", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("base-plan-id")
	return cmd
}

func newBatchUpdateStatesCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
	)

	cmd := &cobra.Command{
		Use:   "batch-update-states",
		Short: "Batch update offer states",
		Long:  "Batch update offer states. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.OffersPath(pkg, productID, basePlanID)+":batchUpdateStates", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	_ = cmd.MarkFlagRequired("product-id")
	_ = cmd.MarkFlagRequired("base-plan-id")
	return cmd
}
