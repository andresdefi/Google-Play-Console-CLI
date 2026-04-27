package auth

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	internalauth "github.com/andresdefi/gpc/internal/auth"
	"github.com/andresdefi/gpc/internal/config"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

const (
	defaultServiceAccountName = "gpc-service-account"
)

var (
	requiredServices = []string{
		"androidpublisher.googleapis.com",
		"playdeveloperreporting.googleapis.com",
		"pubsub.googleapis.com",
	}
	requiredRoles = []string{
		"roles/pubsub.editor",
		"roles/monitoring.viewer",
	}
)

type setupOptions struct {
	projectID          string
	createProjectID    string
	serviceAccountName string
	keyFilePath        string
	nonInteractive     bool
}

type gcloudRunner interface {
	Run(args ...string) (string, error)
}

type execGcloudRunner struct{}

func (execGcloudRunner) Run(args ...string) (string, error) {
	cmd := exec.Command("gcloud", args...) // #nosec G204 - arguments are passed directly without shell expansion.
	out, err := cmd.CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(out)), fmt.Errorf("gcloud %s failed: %w\n%s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

func newSetupCmd() *cobra.Command {
	opts := setupOptions{
		serviceAccountName: defaultServiceAccountName,
	}

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Set up Google Cloud service account authentication",
		Long: `Guide service account setup for the Google Play Developer API.

The wizard uses the gcloud CLI to select or create a Google Cloud project, enable
required APIs, create a service account, grant required IAM roles, generate a JSON
key file, and authenticate gpc with that key.`,
		Example: `  gpc auth setup
  gpc auth setup --project my-play-project --key-file ~/.gpc/play-api.json --non-interactive
  gpc auth setup --create-project my-new-play-project --non-interactive`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runSetup(cmd, execGcloudRunner{}, opts); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.projectID, "project", "", "Google Cloud project ID to use")
	cmd.Flags().StringVar(&opts.createProjectID, "create-project", "", "Google Cloud project ID to create and use")
	cmd.Flags().StringVar(&opts.serviceAccountName, "service-account", defaultServiceAccountName, "Service account name")
	cmd.Flags().StringVar(&opts.keyFilePath, "key-file", "", "Path for the generated service account JSON key")
	cmd.Flags().BoolVar(&opts.nonInteractive, "non-interactive", false, "Disable prompts and require decisions via flags")
	return cmd
}

func runSetup(cmd *cobra.Command, runner gcloudRunner, opts setupOptions) error {
	reader := bufio.NewReader(cmd.InOrStdin())

	if err := ensureGcloudReady(runner); err != nil {
		return exitcode.ConfigError("%v", err)
	}

	projectID, err := selectProject(runner, reader, opts)
	if err != nil {
		return exitcode.ConfigError("%v", err)
	}

	if opts.createProjectID != "" {
		if err := ensureProject(runner, opts.createProjectID); err != nil {
			return exitcode.ConfigError("%v", err)
		}
	}

	if err := enableRequiredServices(runner, projectID); err != nil {
		return exitcode.ConfigError("%v", err)
	}

	serviceAccountName := strings.TrimSpace(opts.serviceAccountName)
	if serviceAccountName == "" {
		serviceAccountName = defaultServiceAccountName
	}
	serviceAccountEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", serviceAccountName, projectID)
	if err := ensureServiceAccount(runner, projectID, serviceAccountName, serviceAccountEmail); err != nil {
		return exitcode.ConfigError("%v", err)
	}

	if err := ensureIAMRoles(runner, projectID, serviceAccountEmail); err != nil {
		return exitcode.ConfigError("%v", err)
	}

	keyFilePath, err := resolveKeyFilePath(opts.keyFilePath, serviceAccountName)
	if err != nil {
		return exitcode.ConfigError("%v", err)
	}
	if err := ensureKeyFile(runner, projectID, serviceAccountEmail, keyFilePath, opts.nonInteractive, reader); err != nil {
		return exitcode.ConfigError("%v", err)
	}

	info, err := internalauth.Login(keyFilePath)
	if err != nil {
		return exitcode.AuthError("could not authenticate with generated key file: %v", err)
	}

	output.Success(fmt.Sprintf("Authenticated as %s (project: %s)", info.ClientEmail, info.ProjectID))
	printPlayConsoleReminder(serviceAccountEmail)
	return nil
}

func ensureGcloudReady(runner gcloudRunner) error {
	if _, err := exec.LookPath("gcloud"); err != nil {
		return errors.New("gcloud CLI is not installed or not on PATH. Install it from https://cloud.google.com/sdk/docs/install, then run 'gcloud auth login'")
	}
	step("Checking gcloud authentication")
	account, err := runner.Run("auth", "list", "--filter=status:ACTIVE", "--format=value(account)")
	if err != nil {
		return err
	}
	if strings.TrimSpace(account) == "" {
		return errors.New("gcloud is not authenticated. Run 'gcloud auth login' and retry")
	}
	done("Authenticated with gcloud as " + firstLine(account))
	return nil
}

func selectProject(runner gcloudRunner, reader *bufio.Reader, opts setupOptions) (string, error) {
	if opts.createProjectID != "" {
		return strings.TrimSpace(opts.createProjectID), nil
	}
	if strings.TrimSpace(opts.projectID) != "" {
		return strings.TrimSpace(opts.projectID), nil
	}
	if opts.nonInteractive {
		return "", errors.New("--project or --create-project is required with --non-interactive")
	}

	step("Detecting Google Cloud project")
	active, _ := runner.Run("config", "get-value", "project")
	active = strings.TrimSpace(stripGcloudUnset(active))
	if active != "" {
		if yes, err := confirm(reader, fmt.Sprintf("Use active gcloud project %q? [Y/n]: ", active), true); err != nil {
			return "", err
		} else if yes {
			done("Using project " + active)
			return active, nil
		}
	}

	projects, err := listProjects(runner)
	if err != nil {
		return "", err
	}
	if len(projects) > 0 {
		fmt.Fprintln(os.Stderr, "Available projects:")
		for i, project := range projects {
			fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, project)
		}
		fmt.Fprintln(os.Stderr, "  n. Create a new project")
		choice, err := prompt(reader, "Choose a project number, project ID, or 'n': ")
		if err != nil {
			return "", err
		}
		choice = strings.TrimSpace(choice)
		if strings.EqualFold(choice, "n") || strings.EqualFold(choice, "new") {
			return promptNewProjectID(reader, runner)
		}
		if index, err := strconv.Atoi(choice); err == nil && index >= 1 && index <= len(projects) {
			return projects[index-1], nil
		}
		if choice != "" {
			return choice, nil
		}
	}

	return promptNewProjectID(reader, runner)
}

func promptNewProjectID(reader *bufio.Reader, runner gcloudRunner) (string, error) {
	projectID, err := prompt(reader, "New Google Cloud project ID: ")
	if err != nil {
		return "", err
	}
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return "", errors.New("project ID is required")
	}
	if err := ensureProject(runner, projectID); err != nil {
		return "", err
	}
	return projectID, nil
}

func ensureProject(runner gcloudRunner, projectID string) error {
	step("Checking Google Cloud project " + projectID)
	if _, err := runner.Run("projects", "describe", projectID, "--format=value(projectId)"); err == nil {
		done("Project already exists")
		return nil
	}
	step("Creating Google Cloud project " + projectID)
	if _, err := runner.Run("projects", "create", projectID); err != nil {
		return err
	}
	done("Created project " + projectID)
	return nil
}

func enableRequiredServices(runner gcloudRunner, projectID string) error {
	step("Checking required APIs")
	enabledOut, err := runner.Run("services", "list", "--enabled", "--project", projectID, "--format=value(config.name)")
	if err != nil {
		return err
	}
	enabled := linesSet(enabledOut)
	for _, service := range requiredServices {
		if enabled[service] {
			done(service + " already enabled")
			continue
		}
		step("Enabling " + service)
		if _, err := runner.Run("services", "enable", service, "--project", projectID); err != nil {
			return err
		}
		done(service + " enabled")
	}
	return nil
}

func ensureServiceAccount(runner gcloudRunner, projectID, name, email string) error {
	step("Checking service account " + email)
	if _, err := runner.Run("iam", "service-accounts", "describe", email, "--project", projectID, "--format=value(email)"); err == nil {
		done("Service account already exists")
		return nil
	}
	step("Creating service account " + name)
	if _, err := runner.Run("iam", "service-accounts", "create", name, "--project", projectID, "--display-name", "gpc service account"); err != nil {
		return err
	}
	done("Created service account " + email)
	return nil
}

func ensureIAMRoles(runner gcloudRunner, projectID, email string) error {
	for _, role := range requiredRoles {
		step("Checking IAM role " + role)
		out, err := runner.Run(
			"projects", "get-iam-policy", projectID,
			"--flatten=bindings[].members",
			"--filter=bindings.role:"+role+" AND bindings.members:serviceAccount:"+email,
			"--format=value(bindings.role)",
		)
		if err == nil && strings.Contains(out, role) {
			done(role + " already granted")
			continue
		}
		step("Granting IAM role " + role)
		if _, err := runner.Run(
			"projects", "add-iam-policy-binding", projectID,
			"--member", "serviceAccount:"+email,
			"--role", role,
		); err != nil {
			return err
		}
		done(role + " granted")
	}
	return nil
}

func ensureKeyFile(runner gcloudRunner, projectID, email, keyPath string, nonInteractive bool, reader *bufio.Reader) error {
	if _, err := os.Stat(keyPath); err == nil {
		if nonInteractive {
			done("Using existing key file " + keyPath)
			return nil
		}
		overwrite, err := confirm(reader, fmt.Sprintf("Key file %q already exists. Overwrite? [y/N]: ", keyPath), false)
		if err != nil {
			return err
		}
		if !overwrite {
			done("Using existing key file " + keyPath)
			return nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(keyPath), 0o700); err != nil {
		return fmt.Errorf("could not create key file directory: %w", err)
	}
	step("Generating service account key")
	if _, err := runner.Run(
		"iam", "service-accounts", "keys", "create", keyPath,
		"--iam-account", email,
		"--project", projectID,
	); err != nil {
		return err
	}
	done("Wrote key file " + keyPath)
	return nil
}

func resolveKeyFilePath(path string, serviceAccountName string) (string, error) {
	if strings.TrimSpace(path) == "" {
		dir, err := config.Dir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(dir, serviceAccountName+".json")
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	return filepath.Abs(path)
}

func listProjects(runner gcloudRunner) ([]string, error) {
	step("Listing Google Cloud projects")
	out, err := runner.Run("projects", "list", "--format=value(projectId)")
	if err != nil {
		return nil, err
	}
	var projects []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			projects = append(projects, line)
		}
	}
	return projects, nil
}

func prompt(reader *bufio.Reader, message string) (string, error) {
	fmt.Fprint(os.Stderr, message)
	value, err := reader.ReadString('\n')
	if err != nil && (!errors.Is(err, io.EOF) || value == "") {
		return "", err
	}
	return strings.TrimSpace(value), nil
}

func confirm(reader *bufio.Reader, message string, defaultYes bool) (bool, error) {
	value, err := prompt(reader, message)
	if err != nil {
		return false, err
	}
	if value == "" {
		return defaultYes, nil
	}
	switch strings.ToLower(value[:1]) {
	case "y":
		return true, nil
	case "n":
		return false, nil
	default:
		return false, fmt.Errorf("expected yes or no, got %q", value)
	}
}

func linesSet(value string) map[string]bool {
	result := make(map[string]bool)
	for _, line := range strings.Split(value, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			result[line] = true
		}
	}
	return result
}

func stripGcloudUnset(value string) string {
	value = strings.TrimSpace(value)
	if strings.Contains(value, "unset") {
		return ""
	}
	return value
}

func firstLine(value string) string {
	if idx := strings.IndexByte(value, '\n'); idx >= 0 {
		return strings.TrimSpace(value[:idx])
	}
	return strings.TrimSpace(value)
}

func step(message string) {
	fmt.Fprintf(os.Stderr, "-> %s...\n", message)
}

func done(message string) {
	fmt.Fprintf(os.Stderr, "   %s\n", message)
}

func printPlayConsoleReminder(email string) {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Next step: invite this service account in Google Play Console > Users and permissions:")
	fmt.Fprintf(os.Stderr, "  %s\n", email)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Grant the app permissions gpc needs:")
	fmt.Fprintln(os.Stderr, "  - View app information and download bulk reports (read-only)")
	fmt.Fprintln(os.Stderr, "  - View financial data, orders, and cancellation survey responses")
	fmt.Fprintln(os.Stderr, "  - Manage orders and subscriptions")
}
