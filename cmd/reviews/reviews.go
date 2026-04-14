package reviews

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
		Use:     "reviews",
		Aliases: []string{"review"},
		Short:   "Manage reviews",
		Long:    "List, get, or reply to user reviews.",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newReplyCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List reviews",
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
			resp, err := client.Get(api.ReviewsPath(pkg), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var result struct {
					Reviews []struct {
						ReviewID string `json:"reviewId"`
						AuthorName string `json:"authorName"`
						Comments []struct {
							UserComment struct {
								StarRating int    `json:"starRating"`
								Text       string `json:"text"`
							} `json:"userComment"`
						} `json:"comments"`
					} `json:"reviews"`
				}
				if err := json.Unmarshal(data.(json.RawMessage), &result); err == nil {
					t := output.NewTable(w, "Review ID", "Author", "Stars", "Comment")
					for _, r := range result.Reviews {
						comment := ""
						stars := 0
						if len(r.Comments) > 0 {
							stars = r.Comments[0].UserComment.StarRating
							comment = r.Comments[0].UserComment.Text
						}
						if len(comment) > 50 {
							comment = comment[:50] + "..."
						}
						t.AppendRow([]any{r.ReviewID, r.AuthorName, stars, comment})
					}
					t.Render()
				} else {
					fmt.Fprintln(w, string(data.(json.RawMessage)))
				}
			})
			return nil
		},
	}
}

func newGetCmd() *cobra.Command {
	var reviewID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a review",
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
			resp, err := client.Get(api.ReviewPath(pkg, reviewID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&reviewID, "review-id", "", "Review ID (required)")
	cmd.MarkFlagRequired("review-id")
	return cmd
}

func newReplyCmd() *cobra.Command {
	var (
		reviewID string
		text     string
	)

	cmd := &cobra.Command{
		Use:   "reply",
		Short: "Reply to a review",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			body := map[string]string{"replyText": text}

			client := api.NewClient(token)
			_, err = client.Post(api.ReviewPath(pkg, reviewID)+":reply", body)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Reply posted to review %s", reviewID))
			return nil
		},
	}
	cmd.Flags().StringVar(&reviewID, "review-id", "", "Review ID (required)")
	cmd.Flags().StringVar(&text, "text", "", "Reply text (required)")
	cmd.MarkFlagRequired("review-id")
	cmd.MarkFlagRequired("text")
	return cmd
}
