package doctor

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/andresdefi/gpc/internal/auth"
	"github.com/andresdefi/gpc/internal/config"
	"github.com/andresdefi/gpc/internal/version"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check your gpc setup",
		Long: `Validate that gpc is configured correctly.

Checks:
  - CLI version and Go runtime
  - Config file exists and is readable
  - Service account credentials are stored
  - OAuth2 token can be obtained
  - Google Play API is reachable`,
		RunE: func(cmd *cobra.Command, args []string) error {
			passed := 0
			failed := 0

			check := func(name string, fn func() (string, error)) {
				detail, err := fn()
				if err != nil {
					fmt.Fprintf(os.Stderr, "  x %s: %v\n", name, err)
					failed++
				} else {
					fmt.Fprintf(os.Stderr, "  + %s: %s\n", name, detail)
					passed++
				}
			}

			fmt.Fprintln(os.Stderr, "Running diagnostics...")
			fmt.Fprintln(os.Stderr)

			check("Version", func() (string, error) {
				return version.String(), nil
			})

			check("Go runtime", func() (string, error) {
				return fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH), nil
			})

			check("Config file", func() (string, error) {
				p, err := config.Path()
				if err != nil {
					return "", err
				}
				if _, err := os.Stat(p); err != nil {
					return "", fmt.Errorf("not found at %s", p)
				}
				return p, nil
			})

			check("Config readable", func() (string, error) {
				cfg, err := config.Load()
				if err != nil {
					return "", err
				}
				if cfg.PackageName != "" {
					return fmt.Sprintf("default package: %s", cfg.PackageName), nil
				}
				return "loaded (no default package set)", nil
			})

			check("Credentials", func() (string, error) {
				info, err := auth.GetStatus()
				if err != nil {
					return "", fmt.Errorf("not configured - run 'gpc auth login'")
				}
				return auth.MaskEmail(info.ClientEmail), nil
			})

			check("OAuth2 token", func() (string, error) {
				_, err := auth.GetToken()
				if err != nil {
					return "", fmt.Errorf("could not obtain token: %v", err)
				}
				return "valid", nil
			})

			check("API reachable", func() (string, error) {
				client := &http.Client{Timeout: 10 * time.Second}
				start := time.Now()
				resp, err := client.Get("https://androidpublisher.googleapis.com/$discovery/rest?version=v3")
				elapsed := time.Since(start)
				if err != nil {
					return "", fmt.Errorf("cannot reach Google Play API: %v", err)
				}
				_ = resp.Body.Close()
				if resp.StatusCode != 200 {
					return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
				}
				return fmt.Sprintf("ok (%dms)", elapsed.Milliseconds()), nil
			})

			check("Environment", func() (string, error) {
				var envs []string
				if v := os.Getenv("GPC_PACKAGE"); v != "" {
					envs = append(envs, "GPC_PACKAGE="+v)
				}
				if v := os.Getenv("GPC_KEY_FILE"); v != "" {
					envs = append(envs, "GPC_KEY_FILE=(set)")
				}
				if v := os.Getenv("GPC_OUTPUT"); v != "" {
					envs = append(envs, "GPC_OUTPUT="+v)
				}
				if _, ok := os.LookupEnv("NO_COLOR"); ok {
					envs = append(envs, "NO_COLOR=(set)")
				}
				if len(envs) == 0 {
					return "no env overrides set", nil
				}
				result := ""
				for _, e := range envs {
					result += e + " "
				}
				return result, nil
			})

			fmt.Fprintln(os.Stderr)
			fmt.Fprintf(os.Stderr, "%d passed, %d failed\n", passed, failed)

			if failed > 0 {
				return fmt.Errorf("%d check(s) failed", failed)
			}
			return nil
		},
	}
}
