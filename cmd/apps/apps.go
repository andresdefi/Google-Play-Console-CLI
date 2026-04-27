package apps

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/andresdefi/gpc/internal/api"
	"github.com/andresdefi/gpc/internal/cmdutil"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apps",
		Aliases: []string{"app"},
		Short:   "Manage apps",
	}

	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newStatusCmd())
	return cmd
}

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get app details",
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
			// Use edits to get app details.
			edit, err := client.CreateEdit(pkg)
			if err != nil {
				return exitcode.APIErrorExit("could not create edit: %v", err)
			}
			defer func() { _ = client.DeleteEdit(pkg, edit.ID) }()

			resp, err := client.Get(api.DetailsPath(pkg, edit.ID), nil)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, json.RawMessage(resp), func(w io.Writer, data any) {
				var d struct {
					DefaultLanguage string `json:"defaultLanguage"`
					ContactEmail    string `json:"contactEmail"`
					ContactPhone    string `json:"contactPhone"`
					ContactWebsite  string `json:"contactWebsite"`
				}
				if err := json.Unmarshal(data.(json.RawMessage), &d); err == nil {
					t := output.NewTable(w, "Field", "Value")
					t.AppendRow([]any{"Package", pkg})
					t.AppendRow([]any{"Default Language", d.DefaultLanguage})
					t.AppendRow([]any{"Contact Email", d.ContactEmail})
					t.AppendRow([]any{"Contact Phone", d.ContactPhone})
					t.AppendRow([]any{"Contact Website", d.ContactWebsite})
					t.Render()
				} else {
					_, _ = fmt.Fprintln(w, string(data.(json.RawMessage)))
				}
			})
			return nil
		},
	}
}

type appStatusCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

type appStatusReport struct {
	Package      string           `json:"package"`
	Checks       []appStatusCheck `json:"checks"`
	ManualChecks []string         `json:"manualChecks"`
	Ready        bool             `json:"ready"`
	Note         string           `json:"note"`
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check app API readiness",
		Long: `Check whether gpc can access the app and query common publishing resources.

Some first-release setup requirements are only visible in Play Console, so this command
also prints a manual checklist for store listing, content rating, and data safety setup.`,
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
			report, err := buildStatusReport(client, pkg)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, report, func(w io.Writer, data any) {
				r := data.(appStatusReport)
				t := output.NewTable(w, "Check", "Status", "Detail")
				for _, check := range r.Checks {
					t.AppendRow([]any{check.Name, check.Status, check.Detail})
				}
				for _, check := range r.ManualChecks {
					t.AppendRow([]any{check, "manual", "Complete or verify in Play Console"})
				}
				t.AppendRow([]any{"Overall", readinessStatus(r.Ready), r.Note})
				t.Render()
			})
			return nil
		},
	}
}

func buildStatusReport(client *api.Client, pkg string) (appStatusReport, error) {
	report := appStatusReport{
		Package: pkg,
		ManualChecks: []string{
			"Store listing completed",
			"Content rating completed",
			"Data safety form completed",
			"App access, ads, and target audience forms completed",
			"Service account invited in Play Console > Users and permissions",
		},
		Ready: true,
		Note:  "API checks passed; complete the manual Play Console checks before first upload or release.",
	}

	edit, err := client.CreateEdit(pkg)
	if err != nil {
		return report, fmt.Errorf("could not create edit: %w", err)
	}
	defer func() { _ = client.DeleteEdit(pkg, edit.ID) }()

	report.addCheck("App access", true, "Created edit "+edit.ID)

	defaultLanguage := ""
	detailsResp, err := client.Get(api.DetailsPath(pkg, edit.ID), nil)
	if err != nil {
		report.addCheck("App details", false, err.Error())
	} else {
		var details struct {
			DefaultLanguage string `json:"defaultLanguage"`
			ContactEmail    string `json:"contactEmail"`
			ContactWebsite  string `json:"contactWebsite"`
		}
		if err := json.Unmarshal(detailsResp, &details); err != nil {
			report.addCheck("App details", false, "Could not parse response: "+err.Error())
		} else {
			defaultLanguage = details.DefaultLanguage
			missing := missingDetailFields(details.DefaultLanguage, details.ContactEmail, details.ContactWebsite)
			if len(missing) == 0 {
				report.addCheck("App details", true, "Default language "+details.DefaultLanguage)
			} else {
				report.addCheck("App details", false, "Missing "+strings.Join(missing, ", "))
			}
		}
	}

	if defaultLanguage != "" {
		if _, err := client.Get(api.ListingPath(pkg, edit.ID, defaultLanguage), nil); err != nil {
			report.addCheck("Default store listing", false, err.Error())
		} else {
			report.addCheck("Default store listing", true, defaultLanguage)
		}
	} else {
		report.addCheck("Default store listing", false, "Skipped because default language is missing")
	}

	tracksResp, err := client.Get(api.TracksPath(pkg, edit.ID), nil)
	if err != nil {
		report.addCheck("Tracks", false, err.Error())
	} else {
		var tracks struct {
			Tracks []struct {
				Track    string `json:"track"`
				Releases []struct {
					Name         string   `json:"name"`
					Status       string   `json:"status"`
					VersionCodes []string `json:"versionCodes"`
				} `json:"releases"`
			} `json:"tracks"`
		}
		if err := json.Unmarshal(tracksResp, &tracks); err != nil {
			report.addCheck("Tracks", false, "Could not parse response: "+err.Error())
		} else if len(tracks.Tracks) == 0 {
			report.addCheck("Tracks", false, "No tracks returned")
		} else {
			report.addCheck("Tracks", true, fmt.Sprintf("%d track(s) returned", len(tracks.Tracks)))
		}
	}

	for _, check := range report.Checks {
		if check.Status != "ok" {
			report.Ready = false
			report.Note = "Some API checks need attention; Play Console may still require manual setup steps."
			break
		}
	}

	return report, nil
}

func (r *appStatusReport) addCheck(name string, ok bool, detail string) {
	status := "ok"
	if !ok {
		status = "needs attention"
	}
	r.Checks = append(r.Checks, appStatusCheck{Name: name, Status: status, Detail: detail})
}

func missingDetailFields(defaultLanguage, contactEmail, contactWebsite string) []string {
	var missing []string
	if strings.TrimSpace(defaultLanguage) == "" {
		missing = append(missing, "defaultLanguage")
	}
	if strings.TrimSpace(contactEmail) == "" {
		missing = append(missing, "contactEmail")
	}
	if strings.TrimSpace(contactWebsite) == "" {
		missing = append(missing, "contactWebsite")
	}
	return missing
}

func readinessStatus(ready bool) string {
	if ready {
		return "ok"
	}
	return "needs attention"
}
