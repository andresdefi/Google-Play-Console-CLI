package releases

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/andresdefi/gpc/internal/spinner"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "releases",
		Aliases: []string{"release"},
		Short:   "Manage releases",
		Long:    "List releases, deploy artifacts, promote between tracks, and manage rollouts.",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newDeployCmd())
	cmd.AddCommand(newPromoteCmd())
	cmd.AddCommand(newRolloutCmd())
	cmd.AddCommand(newHaltCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	var trackName string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List releases for a track",
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

			resp, err := client.Get(api.TrackPath(pkg, edit.ID, trackName), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var t struct {
					Track    string `json:"track"`
					Releases []struct {
						Name         string   `json:"name"`
						Status       string   `json:"status"`
						VersionCodes []string `json:"versionCodes"`
						UserFraction float64  `json:"userFraction"`
						ReleaseNotes []any    `json:"releaseNotes"`
					} `json:"releases"`
				}
				raw := data.(json.RawMessage)
				if err := json.Unmarshal(raw, &t); err == nil {
					tbl := output.NewTable(w, "Name", "Status", "Version Codes", "Rollout %")
					for _, r := range t.Releases {
						pct := "-"
						if r.UserFraction > 0 {
							pct = fmt.Sprintf("%.0f%%", r.UserFraction*100)
						}
						tbl.AppendRow([]any{r.Name, r.Status, strings.Join(r.VersionCodes, ", "), pct})
					}
					tbl.Render()
				} else {
					_, _ = fmt.Fprintln(w, string(raw))
				}
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&trackName, "track", "production", "Track name")
	return cmd
}

func newDeployCmd() *cobra.Command {
	var (
		trackName   string
		rollout     float64
		releaseName string
		notes       string
		notesLang   string
	)

	cmd := &cobra.Command{
		Use:   "deploy <file>",
		Short: "Deploy an APK or AAB to a track",
		Long: `Deploy an APK or AAB file to a release track in one step.

This convenience command handles the full edit flow:
  1. Create an edit session
  2. Upload the APK/AAB
  3. Assign to the specified track
  4. Commit the edit`,
		Example: `  gpc releases deploy app-release.aab --track production
  gpc releases deploy app.aab --track beta --rollout 0.1
  gpc releases deploy app.apk --track internal --release-name "v1.2.3" --notes "Bug fixes"`,
		Args: cobra.ExactArgs(1),
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

			// Validate file exists.
			if _, err := os.Stat(filePath); err != nil {
				return exitcode.NewExitError(exitcode.Error, "file not found: %s", filePath)
			}

			// Detect file type.
			ext := strings.ToLower(filepath.Ext(filePath))
			isBundle := ext == ".aab"
			isAPK := ext == ".apk"
			if !isBundle && !isAPK {
				return exitcode.NewExitError(exitcode.Error, "unsupported file type: %s (expected .apk or .aab)", ext)
			}

			client := api.NewClient(token)
			artifactType := "bundle"
			if isAPK {
				artifactType = "APK"
			}

			sp := spinner.New(fmt.Sprintf("Uploading %s", artifactType))
			sp.Start()

			_, err = client.WithEdit(pkg, func(editID string) error {
				// Step 1: Upload the artifact.
				var uploadPath string
				if isBundle {
					uploadPath = api.BundlesPath(pkg, editID)
				} else {
					uploadPath = api.APKsPath(pkg, editID)
				}

				uploadResp, err := client.Upload(uploadPath, filePath, "application/octet-stream")
				if err != nil {
					sp.Stop("Upload failed")
					return fmt.Errorf("upload failed: %w", err)
				}

				// Extract version code from upload response.
				var uploaded struct {
					VersionCode int `json:"versionCode"`
				}
				if err := json.Unmarshal(uploadResp, &uploaded); err != nil {
					sp.Stop("Upload failed")
					return fmt.Errorf("could not parse upload response: %w", err)
				}
				sp.Stop(fmt.Sprintf("Uploaded version code: %d", uploaded.VersionCode))

				// Step 2: Assign to track.
				sp2 := spinner.New(fmt.Sprintf("Assigning to %s track", trackName))
				sp2.Start()

				release := map[string]any{
					"versionCodes": []int{uploaded.VersionCode},
					"status":       "completed",
				}

				if rollout > 0 && rollout < 1 {
					release["status"] = "inProgress"
					release["userFraction"] = rollout
				}

				if releaseName != "" {
					release["name"] = releaseName
				}

				if notes != "" {
					lang := notesLang
					if lang == "" {
						lang = "en-US"
					}
					release["releaseNotes"] = []map[string]string{
						{"language": lang, "text": notes},
					}
				}

				trackBody := map[string]any{
					"track":    trackName,
					"releases": []any{release},
				}

				_, err = client.Put(api.TrackPath(pkg, editID, trackName), trackBody)
				if err != nil {
					sp2.Stop("Failed")
					return fmt.Errorf("could not assign to track: %w", err)
				}
				sp2.Stop("Assigned")

				return nil
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Successfully deployed %s to %s", filepath.Base(filePath), trackName))
			return nil
		},
	}

	cmd.Flags().StringVar(&trackName, "track", "internal", "Target track (internal, alpha, beta, production, or custom)")
	cmd.Flags().Float64Var(&rollout, "rollout", 0, "Staged rollout fraction (0.0-1.0, 0 = full rollout)")
	cmd.Flags().StringVar(&releaseName, "release-name", "", "Release name (e.g. v1.2.3)")
	cmd.Flags().StringVar(&notes, "notes", "", "Release notes text")
	cmd.Flags().StringVar(&notesLang, "notes-lang", "en-US", "Release notes language code")
	return cmd
}

func newPromoteCmd() *cobra.Command {
	var (
		fromTrack string
		toTrack   string
	)

	cmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote a release between tracks",
		Example: `  gpc releases promote --from beta --to production
  gpc releases promote --from internal --to alpha`,
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
				// Get the source track.
				resp, err := client.Get(api.TrackPath(pkg, editID, fromTrack), nil)
				if err != nil {
					return fmt.Errorf("could not get source track %s: %w", fromTrack, err)
				}

				var sourceTrack struct {
					Releases []json.RawMessage `json:"releases"`
				}
				if err := json.Unmarshal(resp, &sourceTrack); err != nil {
					return fmt.Errorf("could not parse source track: %w", err)
				}

				if len(sourceTrack.Releases) == 0 {
					return fmt.Errorf("no releases found on %s track", fromTrack)
				}

				// Apply the latest release to the target track.
				trackBody := map[string]any{
					"track":    toTrack,
					"releases": sourceTrack.Releases,
				}

				_, err = client.Put(api.TrackPath(pkg, editID, toTrack), trackBody)
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Promoted from %s to %s", fromTrack, toTrack))
			return nil
		},
	}

	cmd.Flags().StringVar(&fromTrack, "from", "", "Source track (required)")
	cmd.Flags().StringVar(&toTrack, "to", "", "Target track (required)")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("to")
	return cmd
}

func newRolloutCmd() *cobra.Command {
	var (
		trackName string
		fraction  float64
	)

	cmd := &cobra.Command{
		Use:   "rollout",
		Short: "Update staged rollout fraction",
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
				resp, err := client.Get(api.TrackPath(pkg, editID, trackName), nil)
				if err != nil {
					return err
				}

				var t map[string]any
				if err := json.Unmarshal(resp, &t); err != nil {
					return err
				}

				releases, ok := t["releases"].([]any)
				if !ok || len(releases) == 0 {
					return fmt.Errorf("no releases found on %s track", trackName)
				}

				// Update the latest release's rollout fraction.
				latest := releases[0].(map[string]any)
				latest["userFraction"] = fraction
				if fraction >= 1.0 {
					latest["status"] = "completed"
					delete(latest, "userFraction")
				} else {
					latest["status"] = "inProgress"
				}

				_, err = client.Put(api.TrackPath(pkg, editID, trackName), t)
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			if fraction >= 1.0 {
				output.Success(fmt.Sprintf("Rollout completed on %s", trackName))
			} else {
				output.Success(fmt.Sprintf("Rollout updated to %.0f%% on %s", fraction*100, trackName))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&trackName, "track", "production", "Track name")
	cmd.Flags().Float64Var(&fraction, "fraction", 0, "Rollout fraction 0.0-1.0 (required)")
	_ = cmd.MarkFlagRequired("fraction")
	return cmd
}

func newHaltCmd() *cobra.Command {
	var trackName string

	cmd := &cobra.Command{
		Use:   "halt",
		Short: "Halt a staged rollout",
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
				resp, err := client.Get(api.TrackPath(pkg, editID, trackName), nil)
				if err != nil {
					return err
				}

				var t map[string]any
				if err := json.Unmarshal(resp, &t); err != nil {
					return err
				}

				releases, ok := t["releases"].([]any)
				if !ok || len(releases) == 0 {
					return fmt.Errorf("no releases found on %s track", trackName)
				}

				latest := releases[0].(map[string]any)
				latest["status"] = "halted"

				_, err = client.Put(api.TrackPath(pkg, editID, trackName), t)
				return err
			})
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			output.Success(fmt.Sprintf("Rollout halted on %s", trackName))
			return nil
		},
	}

	cmd.Flags().StringVar(&trackName, "track", "production", "Track name")
	return cmd
}
