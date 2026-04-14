package listings

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
		Use:     "listings",
		Aliases: []string{"listing"},
		Short:   "Manage store listings",
		Long:    "List, get, update, or delete store listings for an app.",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newDeleteAllCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all store listings",
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
			resp, err := withTempEdit(client, pkg, func(editID string) (json.RawMessage, error) {
				return client.Get(api.ListingsPath(pkg, editID), nil)
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, func(w io.Writer, data any) {
				var result struct {
					Listings []struct {
						Language         string `json:"language"`
						Title            string `json:"title"`
						ShortDescription string `json:"shortDescription"`
					} `json:"listings"`
				}
				if err := json.Unmarshal(data.(json.RawMessage), &result); err == nil {
					t := output.NewTable(w, "Language", "Title", "Short Description")
					for _, l := range result.Listings {
						t.AppendRow([]any{l.Language, l.Title, l.ShortDescription})
					}
					t.Render()
				} else {
					_, _ = fmt.Fprintln(w, string(data.(json.RawMessage)))
				}
			})
			return nil
		},
	}
}

func newGetCmd() *cobra.Command {
	var language string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a store listing by language",
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
			resp, err := withTempEdit(client, pkg, func(editID string) (json.RawMessage, error) {
				return client.Get(api.ListingPath(pkg, editID, language), nil)
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&language, "language", "", "Language code (required)")
	_ = cmd.MarkFlagRequired("language")
	return cmd
}

func newUpdateCmd() *cobra.Command {
	var language string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a store listing",
		Long:  "Update a store listing for a language. Pass the listing data as JSON via stdin.",
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
			_, err = client.WithEdit(pkg, func(editID string) error {
				var body json.RawMessage
				if err := json.NewDecoder(cmd.InOrStdin()).Decode(&body); err != nil {
					return fmt.Errorf("could not read listing data from stdin: %w", err)
				}
				_, err := client.Put(api.ListingPath(pkg, editID, language), body)
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Listing for %s updated and committed", language))
			return nil
		},
	}
	cmd.Flags().StringVar(&language, "language", "", "Language code (required)")
	_ = cmd.MarkFlagRequired("language")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	var language string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a store listing by language",
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
			_, err = client.WithEdit(pkg, func(editID string) error {
				return client.Delete(api.ListingPath(pkg, editID, language))
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Listing for %s deleted and committed", language))
			return nil
		},
	}
	cmd.Flags().StringVar(&language, "language", "", "Language code (required)")
	_ = cmd.MarkFlagRequired("language")
	return cmd
}

func newDeleteAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-all",
		Short: "Delete all store listings",
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
			_, err = client.WithEdit(pkg, func(editID string) error {
				return client.Delete(api.ListingsPath(pkg, editID))
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success("All listings deleted and committed")
			return nil
		},
	}
}

// withTempEdit creates a temporary read-only edit, runs the function, then deletes the edit.
func withTempEdit(client *api.Client, pkg string, fn func(editID string) (json.RawMessage, error)) (json.RawMessage, error) {
	edit, err := client.CreateEdit(pkg)
	if err != nil {
		return nil, err
	}
	defer func() { _ = client.DeleteEdit(pkg, edit.ID) }()

	return fn(edit.ID)
}
