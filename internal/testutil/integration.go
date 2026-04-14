package testutil

import (
	"os"
	"testing"
)

// SkipUnlessIntegration skips the test unless GPC_INTEGRATION_TEST is set.
// Use this for tests that hit the real Google Play API.
func SkipUnlessIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("GPC_INTEGRATION_TEST") == "" {
		t.Skip("skipping integration test; set GPC_INTEGRATION_TEST=1 to run")
	}
}

// RequireKeyFile returns the path to the service account key file.
// Skips the test if GPC_KEY_FILE is not set.
func RequireKeyFile(t *testing.T) string {
	t.Helper()
	SkipUnlessIntegration(t)
	path := os.Getenv("GPC_KEY_FILE")
	if path == "" {
		t.Skip("skipping: GPC_KEY_FILE not set")
	}
	return path
}

// RequirePackage returns the test package name.
// Skips the test if GPC_TEST_PACKAGE is not set.
func RequirePackage(t *testing.T) string {
	t.Helper()
	SkipUnlessIntegration(t)
	pkg := os.Getenv("GPC_TEST_PACKAGE")
	if pkg == "" {
		t.Skip("skipping: GPC_TEST_PACKAGE not set")
	}
	return pkg
}
