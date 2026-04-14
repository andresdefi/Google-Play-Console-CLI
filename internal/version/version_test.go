package version

import (
	"strings"
	"testing"
)

func TestString_Defaults(t *testing.T) {
	// With default values.
	s := String()
	if !strings.Contains(s, "gpc") {
		t.Errorf("expected 'gpc' in version string, got %q", s)
	}
	if !strings.Contains(s, "devel") {
		t.Errorf("expected 'devel' in version string, got %q", s)
	}
	if !strings.Contains(s, "unknown") {
		t.Errorf("expected 'unknown' in version string, got %q", s)
	}
}

func TestString_WithValues(t *testing.T) {
	oldVersion, oldCommit, oldDate := Version, Commit, Date
	defer func() {
		Version, Commit, Date = oldVersion, oldCommit, oldDate
	}()

	Version = "1.2.3"
	Commit = "abc123"
	Date = "2025-01-01"

	s := String()
	expected := "gpc 1.2.3 (commit: abc123, built: 2025-01-01)"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}

func TestVersion_DefaultValue(t *testing.T) {
	if Version != "devel" {
		// It may have been overridden by ldflags, so just check it's non-empty.
		if Version == "" {
			t.Error("expected non-empty Version")
		}
	}
}

func TestCommit_DefaultValue(t *testing.T) {
	if Commit != "unknown" {
		if Commit == "" {
			t.Error("expected non-empty Commit")
		}
	}
}

func TestDate_DefaultValue(t *testing.T) {
	if Date != "unknown" {
		if Date == "" {
			t.Error("expected non-empty Date")
		}
	}
}
