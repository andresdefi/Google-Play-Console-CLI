package expansionfiles

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "expansionfiles",
		Aliases: []string{"expansion-file", "obb"},
		Short:   "Manage expansion files",
		Long:    "Get, update, or upload expansion files (OBBs) for an app.",
	}

	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newUploadCmd())
	return cmd
}

func newGetCmd() *cobra.Command {
	var (
		apkVersion int
		fileType   string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an expansion file",
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
			edit, err := client.CreateEdit(pkg)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}
			defer func() { _ = client.DeleteEdit(pkg, edit.ID) }()

			resp, err := client.Get(api.ExpansionFilesPath(pkg, edit.ID, apkVersion, fileType), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().IntVar(&apkVersion, "apk-version", 0, "APK version code (required)")
	cmd.Flags().StringVar(&fileType, "type", "", "Expansion file type: main or patch (required)")
	_ = cmd.MarkFlagRequired("apk-version")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func newUpdateCmd() *cobra.Command {
	var (
		apkVersion int
		fileType   string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an expansion file",
		Long:  "Update an expansion file configuration. Pass the data as JSON via stdin.",
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
				return exitcode.ConfigError("could not read expansion file data from stdin: %v", err)
			}

			client := api.NewClient(token)
			_, err = client.WithEdit(pkg, func(editID string) error {
				_, err := client.Put(api.ExpansionFilesPath(pkg, editID, apkVersion, fileType), body)
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Expansion file updated and committed for %s", pkg))
			return nil
		},
	}
	cmd.Flags().IntVar(&apkVersion, "apk-version", 0, "APK version code (required)")
	cmd.Flags().StringVar(&fileType, "type", "", "Expansion file type: main or patch (required)")
	_ = cmd.MarkFlagRequired("apk-version")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func newUploadCmd() *cobra.Command {
	var (
		apkVersion int
		fileType   string
	)

	cmd := &cobra.Command{
		Use:   "upload <file>",
		Short: "Upload an expansion file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			if _, err := os.Stat(filePath); err != nil {
				return exitcode.NewExitError(exitcode.Error, "file not found: %s", filePath)
			}

			client := api.NewClient(token)
			_, err = client.WithEdit(pkg, func(editID string) error {
				_, err := client.Upload(api.ExpansionFilesPath(pkg, editID, apkVersion, fileType), filePath, "application/octet-stream")
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Expansion file uploaded and committed for %s", pkg))
			return nil
		},
	}
	cmd.Flags().IntVar(&apkVersion, "apk-version", 0, "APK version code (required)")
	cmd.Flags().StringVar(&fileType, "type", "", "Expansion file type: main or patch (required)")
	_ = cmd.MarkFlagRequired("apk-version")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}
