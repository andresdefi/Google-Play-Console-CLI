package cmd

import (
	"testing"
)

func findSubcommand(name string) bool {
	for _, c := range rootCmd.Commands() {
		if c.Use == name || c.Name() == name {
			return true
		}
	}
	return false
}

func TestRootCommand_Exists(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}
	if rootCmd.Use != "gpc" {
		t.Errorf("expected Use 'gpc', got %q", rootCmd.Use)
	}
}

func TestRootCommand_HasSubcommands(t *testing.T) {
	if len(rootCmd.Commands()) == 0 {
		t.Fatal("rootCmd should have subcommands")
	}
}

func TestRootCommand_HasPersistentFlags_Package(t *testing.T) {
	f := rootCmd.PersistentFlags().Lookup("package")
	if f == nil {
		t.Fatal("expected persistent flag 'package'")
	}
	if f.Shorthand != "p" {
		t.Errorf("expected shorthand 'p', got %q", f.Shorthand)
	}
}

func TestRootCommand_HasPersistentFlags_Output(t *testing.T) {
	f := rootCmd.PersistentFlags().Lookup("output")
	if f == nil {
		t.Fatal("expected persistent flag 'output'")
	}
	if f.Shorthand != "o" {
		t.Errorf("expected shorthand 'o', got %q", f.Shorthand)
	}
}

func TestAuthCommand_Exists(t *testing.T) {
	if !findSubcommand("auth") {
		t.Error("expected 'auth' subcommand")
	}
}

func TestAuthCommand_HasSubcommands_Login(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "auth" {
			found := false
			for _, sub := range c.Commands() {
				if sub.Name() == "login" {
					found = true
					break
				}
			}
			if !found {
				t.Error("expected 'login' subcommand under 'auth'")
			}
			return
		}
	}
	t.Error("auth command not found")
}

func TestAuthCommand_HasSubcommands_Status(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "auth" {
			found := false
			for _, sub := range c.Commands() {
				if sub.Name() == "status" {
					found = true
					break
				}
			}
			if !found {
				t.Error("expected 'status' subcommand under 'auth'")
			}
			return
		}
	}
	t.Error("auth command not found")
}

func TestAuthCommand_HasSubcommands_Logout(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "auth" {
			found := false
			for _, sub := range c.Commands() {
				if sub.Name() == "logout" {
					found = true
					break
				}
			}
			if !found {
				t.Error("expected 'logout' subcommand under 'auth'")
			}
			return
		}
	}
	t.Error("auth command not found")
}

func TestVersionCommand_Exists(t *testing.T) {
	if !findSubcommand("version") {
		t.Error("expected 'version' subcommand")
	}
}

func TestAppsCommand_Exists(t *testing.T) {
	if !findSubcommand("apps") {
		t.Error("expected 'apps' subcommand")
	}
}

func TestEditsCommand_Exists(t *testing.T) {
	if !findSubcommand("edits") {
		t.Error("expected 'edits' subcommand")
	}
}

func TestTracksCommand_Exists(t *testing.T) {
	if !findSubcommand("tracks") {
		t.Error("expected 'tracks' subcommand")
	}
}

func TestReleasesCommand_Exists(t *testing.T) {
	if !findSubcommand("releases") {
		t.Error("expected 'releases' subcommand")
	}
}

func TestIAPCommand_Exists(t *testing.T) {
	if !findSubcommand("iap") {
		t.Error("expected 'iap' subcommand")
	}
}

func TestSubscriptionsCommand_Exists(t *testing.T) {
	if !findSubcommand("subscriptions") {
		t.Error("expected 'subscriptions' subcommand")
	}
}

func TestReviewsCommand_Exists(t *testing.T) {
	if !findSubcommand("reviews") {
		t.Error("expected 'reviews' subcommand")
	}
}

func TestOrdersCommand_Exists(t *testing.T) {
	if !findSubcommand("orders") {
		t.Error("expected 'orders' subcommand")
	}
}

func TestPurchasesCommand_Exists(t *testing.T) {
	if !findSubcommand("purchases") {
		t.Error("expected 'purchases' subcommand")
	}
}

func TestUsersCommand_Exists(t *testing.T) {
	if !findSubcommand("users") {
		t.Error("expected 'users' subcommand")
	}
}

func TestListingsCommand_Exists(t *testing.T) {
	if !findSubcommand("listings") {
		t.Error("expected 'listings' subcommand")
	}
}

func TestDevicesCommand_Exists(t *testing.T) {
	if !findSubcommand("devices") {
		t.Error("expected 'devices' subcommand")
	}
}

func TestBundlesCommand_Exists(t *testing.T) {
	if !findSubcommand("bundles") {
		t.Error("expected 'bundles' subcommand")
	}
}

func TestAPKsCommand_Exists(t *testing.T) {
	if !findSubcommand("apks") {
		t.Error("expected 'apks' subcommand")
	}
}

func TestImagesCommand_Exists(t *testing.T) {
	if !findSubcommand("images") {
		t.Error("expected 'images' subcommand")
	}
}

func TestDetailsCommand_Exists(t *testing.T) {
	if !findSubcommand("details") {
		t.Error("expected 'details' subcommand")
	}
}

func TestTestersCommand_Exists(t *testing.T) {
	if !findSubcommand("testers") {
		t.Error("expected 'testers' subcommand")
	}
}

func TestGrantsCommand_Exists(t *testing.T) {
	if !findSubcommand("grants") {
		t.Error("expected 'grants' subcommand")
	}
}

func TestAllCommandsRegistered(t *testing.T) {
	expected := []string{
		"auth", "version", "apps", "edits", "tracks", "releases",
		"apks", "bundles", "deobfuscation", "expansionfiles", "countryavailability",
		"iap", "subscriptions", "baseplans", "offers",
		"onetimeproducts", "purchaseoptions", "otpoffers", "pricing",
		"listings", "images", "details", "testers", "reviews", "datasafety",
		"orders", "purchases",
		"users", "grants",
		"devices",
		"apprecovery",
		"externaltransactions",
		"generatedapks", "systemapks", "internalsharing",
	}

	cmds := rootCmd.Commands()
	cmdNames := make(map[string]bool)
	for _, c := range cmds {
		cmdNames[c.Name()] = true
	}

	for _, name := range expected {
		if !cmdNames[name] {
			t.Errorf("expected subcommand %q not found", name)
		}
	}

	// Check total count is at least as many as expected.
	if len(cmds) < len(expected) {
		t.Errorf("expected at least %d subcommands, got %d", len(expected), len(cmds))
	}
}
