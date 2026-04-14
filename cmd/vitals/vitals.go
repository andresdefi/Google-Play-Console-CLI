package vitals

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

// The Play Vitals data is accessed through the Play Developer Reporting API (v1beta1).
// Base URL: https://playdeveloperreporting.googleapis.com/v1beta1
// This is separate from the Android Publisher API but uses the same auth.

const reportingBaseURL = "https://playdeveloperreporting.googleapis.com/v1beta1"

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vitals",
		Short: "App vitals and quality metrics",
		Long: `View crash rates, ANR rates, startup performance, and other quality
metrics from the Play Developer Reporting API.

These commands provide the same data visible in Play Console's "Android Vitals"
section, accessible from your terminal or CI/CD pipelines.`,
	}

	cmd.AddCommand(newOverviewCmd())
	cmd.AddCommand(newCrashesCmd())
	cmd.AddCommand(newAnrsCmd())
	cmd.AddCommand(newStartupCmd())
	cmd.AddCommand(newRenderingCmd())
	cmd.AddCommand(newBatteryCmd())
	cmd.AddCommand(newErrorsCmd())
	return cmd
}

func newOverviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "overview",
		Short: "Show vitals overview for an app",
		Long:  "Display a summary of all vital metrics: crash rate, ANR rate, and excessive wakeups.",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			client := api.NewClientWithHTTP(token, nil, reportingBaseURL)
			_ = client

			// The reporting API requires specific query structure.
			// For now, provide the anomaly rate sets endpoint.
			path := fmt.Sprintf("/apps/%s/anomalyRateTimeSeries", pkg)
			resp, err := fetchReportingAPI(token, path)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, func(w io.Writer, data any) {
				_, _ = fmt.Fprintln(w, "Vitals overview for", pkg)
				_, _ = fmt.Fprintln(w, "(Use 'gpc vitals crashes', 'gpc vitals anrs', etc. for detailed metrics)")
				_, _ = fmt.Fprintln(w)
				printRawJSON(w, data)
			})
			return nil
		},
	}
}

func newCrashesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "crashes",
		Short: "View crash rate metrics",
		Long:  "Display crash rate data including user-perceived crash rate and crash clusters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			path := fmt.Sprintf("/apps/%s/crashRateMetricSet:query", pkg)
			resp, err := queryReportingAPI(token, path, defaultMetricQuery("crashRate"))
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, func(w io.Writer, data any) {
				_, _ = fmt.Fprintln(w, "Crash rate metrics for", pkg)
				_, _ = fmt.Fprintln(w)
				printRawJSON(w, data)
			})
			return nil
		},
	}
}

func newAnrsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "anrs",
		Short: "View ANR (Application Not Responding) rate metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			path := fmt.Sprintf("/apps/%s/anrRateMetricSet:query", pkg)
			resp, err := queryReportingAPI(token, path, defaultMetricQuery("anrRate"))
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, func(w io.Writer, data any) {
				_, _ = fmt.Fprintln(w, "ANR rate metrics for", pkg)
				_, _ = fmt.Fprintln(w)
				printRawJSON(w, data)
			})
			return nil
		},
	}
}

func newStartupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "startup",
		Short: "View app startup time metrics",
		Long:  "Display slow and excessive startup time percentages.",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			path := fmt.Sprintf("/apps/%s/slowStartRateMetricSet:query", pkg)
			resp, err := queryReportingAPI(token, path, defaultMetricQuery("slowStartRate"))
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, func(w io.Writer, data any) {
				_, _ = fmt.Fprintln(w, "Startup time metrics for", pkg)
				_, _ = fmt.Fprintln(w)
				printRawJSON(w, data)
			})
			return nil
		},
	}
}

func newRenderingCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rendering",
		Short: "View slow rendering metrics",
		Long:  "Display slow and frozen frame rates.",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			path := fmt.Sprintf("/apps/%s/slowRenderingRateMetricSet:query", pkg)
			resp, err := queryReportingAPI(token, path, defaultMetricQuery("slowRenderingRate"))
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, func(w io.Writer, data any) {
				_, _ = fmt.Fprintln(w, "Rendering metrics for", pkg)
				_, _ = fmt.Fprintln(w)
				printRawJSON(w, data)
			})
			return nil
		},
	}
}

func newBatteryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "battery",
		Short: "View battery usage metrics",
		Long:  "Display excessive wakeup and stuck wake lock rates.",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			path := fmt.Sprintf("/apps/%s/excessiveWakeupRateMetricSet:query", pkg)
			resp, err := queryReportingAPI(token, path, defaultMetricQuery("excessiveWakeupRate"))
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, func(w io.Writer, data any) {
				_, _ = fmt.Fprintln(w, "Battery metrics for", pkg)
				_, _ = fmt.Fprintln(w)
				printRawJSON(w, data)
			})
			return nil
		},
	}
}

func newErrorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "errors",
		Short: "View error counts and clusters",
		Long:  "Query error reports from the Play Developer Reporting API.",
	}

	cmd.AddCommand(newErrorCountsCmd())
	cmd.AddCommand(newErrorIssuesCmd())
	return cmd
}

func newErrorCountsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "counts",
		Short: "View error counts over time",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			path := fmt.Sprintf("/apps/%s/errorCountMetricSet:query", pkg)
			resp, err := queryReportingAPI(token, path, defaultMetricQuery("errorCount"))
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, nil)
			return nil
		},
	}
}

func newErrorIssuesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "issues",
		Short: "Search error issues (crash/ANR clusters)",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkg, err := cmdutil.ResolvePackage(cmd)
			if err != nil {
				return exitcode.ConfigError("%v", err)
			}
			token, err := cmdutil.RequireAuth()
			if err != nil {
				return exitcode.AuthError("%v", err)
			}

			path := fmt.Sprintf("/apps/%s/errorIssues:search", pkg)
			resp, err := fetchReportingAPI(token, path)
			if err != nil {
				return exitcode.APIErrorExit("%v", err)
			}

			format := cmdutil.GetOutputFormat(cmd)
			output.Print(format, resp, nil)
			return nil
		},
	}
}

// --- helpers ---

func fetchReportingAPI(token, path string) (json.RawMessage, error) {
	client := api.NewClientWithHTTP(token, nil, reportingBaseURL)
	return client.Get(path, nil)
}

func queryReportingAPI(token, path string, body any) (json.RawMessage, error) {
	client := api.NewClientWithHTTP(token, nil, reportingBaseURL)
	return client.Post(path, body)
}

func defaultMetricQuery(metric string) map[string]any {
	return map[string]any{
		"metrics": []string{metric},
		"timelineSpec": map[string]any{
			"aggregationPeriod": "DAILY",
		},
	}
}

func printRawJSON(w io.Writer, data any) {
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		_, _ = fmt.Fprintln(w, data)
		return
	}
	_, _ = fmt.Fprintln(w, string(raw))
}
