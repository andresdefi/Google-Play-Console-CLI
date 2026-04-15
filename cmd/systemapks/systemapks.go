package systemapks

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
		Use:     "systemapks",
		Aliases: []string{"system-apk"},
		Short:   "Manage system APK variants",
		Long:    "List, get, create, or download system APK variants for an app.",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newDownloadCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	var versionCode int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List system APK variants",
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
			resp, err := client.Get(api.SystemAPKVariantsPath(pkg, versionCode), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var list struct {
					Variants []struct {
						VariantID  int             `json:"variantId"`
						DeviceSpec json.RawMessage `json:"deviceSpec"`
					} `json:"variants"`
				}
				raw := data.(json.RawMessage)
				if err := json.Unmarshal(raw, &list); err == nil && len(list.Variants) > 0 {
					t := output.NewTable(w, "Variant ID")
					for _, v := range list.Variants {
						t.AppendRow([]any{v.VariantID})
					}
					t.Render()
				} else {
					_, _ = fmt.Fprintln(w, string(raw))
				}
			})
			return nil
		},
	}
	cmd.Flags().IntVar(&versionCode, "version-code", 0, "Version code (required)")
	_ = cmd.MarkFlagRequired("version-code")
	return cmd
}

func newGetCmd() *cobra.Command {
	var (
		versionCode int
		variantID   string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a system APK variant",
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
			resp, err := client.Get(api.SystemAPKVariantPath(pkg, versionCode, variantID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().IntVar(&versionCode, "version-code", 0, "Version code (required)")
	cmd.Flags().StringVar(&variantID, "variant-id", "", "Variant ID (required)")
	_ = cmd.MarkFlagRequired("version-code")
	_ = cmd.MarkFlagRequired("variant-id")
	return cmd
}

func newCreateCmd() *cobra.Command {
	var versionCode int

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a system APK variant",
		Long:  "Create a new system APK variant. Pass the variant data as JSON via stdin.",
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
				return exitcode.ConfigError("could not read variant data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.SystemAPKVariantsPath(pkg, versionCode), body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().IntVar(&versionCode, "version-code", 0, "Version code (required)")
	_ = cmd.MarkFlagRequired("version-code")
	return cmd
}

func newDownloadCmd() *cobra.Command {
	var (
		versionCode int
		variantID   string
		outputPath  string
	)

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download a system APK variant",
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
			err = client.DownloadToFile(api.SystemAPKVariantPath(pkg, versionCode, variantID)+":download", outputPath)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Downloaded to %s", outputPath))
			return nil
		},
	}
	cmd.Flags().IntVar(&versionCode, "version-code", 0, "Version code (required)")
	cmd.Flags().StringVar(&variantID, "variant-id", "", "Variant ID (required)")
	cmd.Flags().StringVar(&outputPath, "output", "", "Output file path (required)")
	_ = cmd.MarkFlagRequired("version-code")
	_ = cmd.MarkFlagRequired("variant-id")
	_ = cmd.MarkFlagRequired("output")
	return cmd
}
