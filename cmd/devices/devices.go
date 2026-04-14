package devices

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
		Use:     "devices",
		Aliases: []string{"device"},
		Short:   "Manage device tier configs",
		Long:    "List, get, or create device tier configurations for an app.",
	}

	cmd.AddCommand(newListTierConfigsCmd())
	cmd.AddCommand(newGetTierConfigCmd())
	cmd.AddCommand(newCreateTierConfigCmd())
	return cmd
}

func newListTierConfigsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-tier-configs",
		Short: "List device tier configs",
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
			resp, err := client.Get(api.DeviceTierConfigsPath(pkg), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var list struct {
					DeviceTierConfigs []struct {
						DeviceTierConfigID string `json:"deviceTierConfigId"`
					} `json:"deviceTierConfigs"`
				}
				raw := data.(json.RawMessage)
				if err := json.Unmarshal(raw, &list); err == nil && len(list.DeviceTierConfigs) > 0 {
					t := output.NewTable(w, "Config ID")
					for _, c := range list.DeviceTierConfigs {
						t.AppendRow([]any{c.DeviceTierConfigID})
					}
					t.Render()
				} else {
					_, _ = fmt.Fprintln(w, string(raw))
				}
			})
			return nil
		},
	}
}

func newGetTierConfigCmd() *cobra.Command {
	var configID string

	cmd := &cobra.Command{
		Use:   "get-tier-config",
		Short: "Get a device tier config",
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
			resp, err := client.Get(api.DeviceTierConfigPath(pkg, configID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&configID, "config-id", "", "Device tier config ID (required)")
	_ = cmd.MarkFlagRequired("config-id")
	return cmd
}

func newCreateTierConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-tier-config",
		Short: "Create a device tier config",
		Long:  "Create a new device tier config. Pass the config data as JSON via stdin.",
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
				return exitcode.ConfigError("could not read config data from stdin: %v", err)
			}

			client := api.NewClient(token)
			resp, err := client.Post(api.DeviceTierConfigsPath(pkg), body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
}
