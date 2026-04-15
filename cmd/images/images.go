package images

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
		Use:     "images",
		Aliases: []string{"image", "img"},
		Short:   "Manage store listing images",
		Long:    "List, upload, or delete images for store listings.",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newUploadCmd())
	cmd.AddCommand(newDeleteCmd())
	cmd.AddCommand(newDeleteAllCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	var (
		language  string
		imageType string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List images for a listing",
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
				return client.Get(api.ImagesPath(pkg, editID, language, imageType), nil)
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, func(w io.Writer, data any) {
				var result struct {
					Images []struct {
						ID     string `json:"id"`
						URL    string `json:"url"`
						SHA1   string `json:"sha1"`
						SHA256 string `json:"sha256"`
					} `json:"images"`
				}
				if err := json.Unmarshal(data.(json.RawMessage), &result); err == nil {
					t := output.NewTable(w, "ID", "URL")
					for _, img := range result.Images {
						t.AppendRow([]any{img.ID, img.URL})
					}
					t.Render()
				} else {
					_, _ = fmt.Fprintln(w, string(data.(json.RawMessage)))
				}
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&language, "language", "", "Language code (required)")
	cmd.Flags().StringVar(&imageType, "type", "", "Image type (required)")
	_ = cmd.MarkFlagRequired("language")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func newUploadCmd() *cobra.Command {
	var (
		language  string
		imageType string
	)

	cmd := &cobra.Command{
		Use:   "upload <file>",
		Short: "Upload an image",
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

			client := api.NewClient(token)
			_, err = client.WithEdit(pkg, func(editID string) error {
				_, err := client.Upload(api.ImagesPath(pkg, editID, language, imageType), filePath, "image/*")
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Image uploaded for %s/%s and committed", language, imageType))
			return nil
		},
	}
	cmd.Flags().StringVar(&language, "language", "", "Language code (required)")
	cmd.Flags().StringVar(&imageType, "type", "", "Image type (required)")
	_ = cmd.MarkFlagRequired("language")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func newDeleteCmd() *cobra.Command {
	var (
		language  string
		imageType string
		imageID   string
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a specific image",
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
				return client.Delete(api.ImagePath(pkg, editID, language, imageType, imageID))
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Image %s deleted and committed", imageID))
			return nil
		},
	}
	cmd.Flags().StringVar(&language, "language", "", "Language code (required)")
	cmd.Flags().StringVar(&imageType, "type", "", "Image type (required)")
	cmd.Flags().StringVar(&imageID, "image-id", "", "Image ID (required)")
	_ = cmd.MarkFlagRequired("language")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("image-id")
	return cmd
}

func newDeleteAllCmd() *cobra.Command {
	var (
		language  string
		imageType string
	)

	cmd := &cobra.Command{
		Use:   "delete-all",
		Short: "Delete all images for a type",
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
				return client.Delete(api.ImagesPath(pkg, editID, language, imageType))
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("All %s images for %s deleted and committed", imageType, language))
			return nil
		},
	}
	cmd.Flags().StringVar(&language, "language", "", "Language code (required)")
	cmd.Flags().StringVar(&imageType, "type", "", "Image type (required)")
	_ = cmd.MarkFlagRequired("language")
	_ = cmd.MarkFlagRequired("type")
	return cmd
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
