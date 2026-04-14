package tracks

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

type track struct {
	Track    string          `json:"track"`
	Releases json.RawMessage `json:"releases,omitempty"`
}

type trackList struct {
	Kind   string  `json:"kind"`
	Tracks []track `json:"tracks"`
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tracks",
		Aliases: []string{"track"},
		Short:   "Manage release tracks",
		Long:    "List, get, update, or create release tracks (internal, alpha, beta, production, custom).",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newCreateCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tracks",
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
			editID, err := withTempEdit(client, pkg, func(editID string) (json.RawMessage, error) {
				return client.Get(api.TracksPath(pkg, editID), nil)
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, editID, func(w io.Writer, data any) {
				var tl trackList
				if err := json.Unmarshal(data.(json.RawMessage), &tl); err == nil {
					t := output.NewTable(w, "Track")
					for _, tr := range tl.Tracks {
						t.AppendRow([]any{tr.Track})
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
	var trackName string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a track",
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
				return client.Get(api.TrackPath(pkg, editID, trackName), nil)
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, nil)
			return nil
		},
	}
	cmd.Flags().StringVar(&trackName, "track", "", "Track name (required)")
	cmd.MarkFlagRequired("track")
	return cmd
}

func newUpdateCmd() *cobra.Command {
	var trackName string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a track",
		Long:  "Update a track configuration. Pass the track data as JSON via stdin.",
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
					return fmt.Errorf("could not read track data from stdin: %w", err)
				}
				_, err := client.Put(api.TrackPath(pkg, editID, trackName), body)
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Track %s updated and committed", trackName))
			return nil
		},
	}
	cmd.Flags().StringVar(&trackName, "track", "", "Track name (required)")
	cmd.MarkFlagRequired("track")
	return cmd
}

func newCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a custom track",
		Long:  "Create a new custom track. Pass the track data as JSON via stdin.",
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
					return fmt.Errorf("could not read track data from stdin: %w", err)
				}
				_, err := client.Post(api.TracksPath(pkg, editID), body)
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success("Track created and committed")
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
	defer client.DeleteEdit(pkg, edit.ID)

	return fn(edit.ID)
}
