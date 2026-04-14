package onetimeproducts

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
		Use:     "onetimeproducts",
		Aliases: []string{"one-time-product", "otp"},
		Short:   "Manage one-time products",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newBatchGetCmd())
	cmd.AddCommand(newBatchUpdateCmd())
	cmd.AddCommand(newBatchDeleteCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List one-time products",
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
			resp, err := client.Get(api.OneTimeProductsPath(pkg), nil)
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
	var productID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a one-time product",
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
			resp, err := client.Get(api.OneTimeProductPath(pkg, productID), nil)
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

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a one-time product",
		Long:  "Create a one-time product. Pass the product data as JSON via stdin.",
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
			resp, err := client.Post(api.OneTimeProductsPath(pkg), body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
}

func newUpdateCmd() *cobra.Command {
	var productID string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a one-time product",
		Long:  "Update a one-time product. Pass the product data as JSON via stdin.",
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
			resp, err := client.Patch(api.OneTimeProductPath(pkg, productID), nil, body)
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

func newDeleteCmd() *cobra.Command {
	var productID string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a one-time product",
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
			if err := client.Delete(api.OneTimeProductPath(pkg, productID)); err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("One-time product %s deleted", productID))
			return nil
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "One-time product ID (required)")
	cmd.MarkFlagRequired("product-id")
	return cmd
}

func newBatchGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-get",
		Short: "Batch get one-time products",
		Long:  "Batch get one-time products. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.OneTimeProductsPath(pkg)+":batchGet", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
}

func newBatchUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-update",
		Short: "Batch update one-time products",
		Long:  "Batch update one-time products. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.OneTimeProductsPath(pkg)+":batchUpdate", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
}

func newBatchDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch-delete",
		Short: "Batch delete one-time products",
		Long:  "Batch delete one-time products. Pass the request body as JSON via stdin.",
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
			resp, err := client.Post(api.OneTimeProductsPath(pkg)+":batchDelete", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
}
