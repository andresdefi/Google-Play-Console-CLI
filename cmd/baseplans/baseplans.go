package baseplans

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
		Use:     "baseplans",
		Aliases: []string{"base-plan", "bp"},
		Short:   "Manage subscription base plans",
	}

	cmd.AddCommand(newActivateCmd())
	cmd.AddCommand(newDeactivateCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newMigratePricesCmd())
	cmd.AddCommand(newBatchMigratePricesCmd())
	cmd.AddCommand(newBatchUpdateStatesCmd())
	return cmd
}

func newActivateCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
	)

	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Activate a base plan",
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
			_, err = client.Post(api.BasePlanPath(pkg, productID, basePlanID)+":activate", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Base plan %s activated", basePlanID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	cmd.MarkFlagRequired("product-id")
	cmd.MarkFlagRequired("base-plan-id")
	return cmd
}

func newDeactivateCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
	)

	cmd := &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate a base plan",
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
			_, err = client.Post(api.BasePlanPath(pkg, productID, basePlanID)+":deactivate", nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Base plan %s deactivated", basePlanID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	cmd.MarkFlagRequired("product-id")
	cmd.MarkFlagRequired("base-plan-id")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a base plan",
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
			if err := client.Delete(api.BasePlanPath(pkg, productID, basePlanID)); err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Base plan %s deleted", basePlanID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.Flags().StringVar(&basePlanID, "base-plan-id", "", "Base plan ID (required)")
	cmd.MarkFlagRequired("product-id")
	cmd.MarkFlagRequired("base-plan-id")
	return cmd
}

func newMigratePricesCmd() *cobra.Command {
	var (
		productID  string
		basePlanID string
	)

	cmd := &cobra.Command{
		Use:   "migrate-prices",
		Short: "Migrate prices for a base plan",
		Long:  "Migrate prices for a base plan. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.BasePlanPath(pkg, productID, basePlanID)+":migratePrices", body)
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
	cmd.MarkFlagRequired("product-id")
	cmd.MarkFlagRequired("base-plan-id")
	return cmd
}

func newBatchMigratePricesCmd() *cobra.Command {
	var productID string

	cmd := &cobra.Command{
		Use:   "batch-migrate-prices",
		Short: "Batch migrate prices for base plans",
		Long:  "Batch migrate prices for base plans. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.BasePlansPath(pkg, productID)+":batchMigratePrices", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.MarkFlagRequired("product-id")
	return cmd
}

func newBatchUpdateStatesCmd() *cobra.Command {
	var productID string

	cmd := &cobra.Command{
		Use:   "batch-update-states",
		Short: "Batch update base plan states",
		Long:  "Batch update base plan states. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.BasePlansPath(pkg, productID)+":batchUpdateStates", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "Subscription product ID (required)")
	cmd.MarkFlagRequired("product-id")
	return cmd
}
