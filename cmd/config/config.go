package config

import (
	"fmt"
	"io"

	internalconfig "github.com/andresdefi/gpc/internal/config"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long:  "Get and set gpc configuration values stored in ~/.gpc/config.json.",
	}

	cmd.AddCommand(newSetCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newPathCmd())
	return cmd
}

func newSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long: `Set a configuration value.

Available keys:
  package    Default Android package name (avoids repeating -p flag)`,
		Example: `  gpc config set package com.example.app`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]

			cfg, err := internalconfig.Load()
			if err != nil {
				return exitcode.ConfigError("could not load config: %v", err)
			}

			switch key {
			case "package":
				cfg.PackageName = value
			default:
				return exitcode.ConfigError("unknown config key: %s", key)
			}

			if err := internalconfig.Save(cfg); err != nil {
				return exitcode.ConfigError("could not save config: %v", err)
			}

			output.Success(fmt.Sprintf("Set %s = %s", key, value))
			return nil
		},
	}
}

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			cfg, err := internalconfig.Load()
			if err != nil {
				return exitcode.ConfigError("could not load config: %v", err)
			}

			switch key {
			case "package":
				if cfg.PackageName == "" {
					return exitcode.ConfigError("package is not set")
				}
				fmt.Println(cfg.PackageName)
			case "key_file_path":
				if cfg.KeyFilePath == "" {
					return exitcode.ConfigError("key_file_path is not set")
				}
				fmt.Println(cfg.KeyFilePath)
			default:
				return exitcode.ConfigError("unknown config key: %s", key)
			}
			return nil
		},
	}
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := internalconfig.Load()
			if err != nil {
				return exitcode.ConfigError("could not load config: %v", err)
			}

			format := output.DetectFormat("")
			if f, _ := cmd.Flags().GetString("output"); f != "" {
				format = output.DetectFormat(f)
			}

			entries := []configEntry{
				{Key: "package", Value: cfg.PackageName},
				{Key: "key_file_path", Value: cfg.KeyFilePath},
			}

			output.Print(format, entries, func(w io.Writer, data any) {
				items := data.([]configEntry)
				t := output.NewTable(w, "Key", "Value")
				for _, e := range items {
					val := e.Value
					if val == "" {
						val = "(not set)"
					}
					t.AppendRow([]any{e.Key, val})
				}
				t.Render()
			})
			return nil
		},
	}
}

func newPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the config file path",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := internalconfig.Path()
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			fmt.Println(p)
			return nil
		},
	}
}

type configEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
