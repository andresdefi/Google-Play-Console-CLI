package cmdutil

import (
	"testing"

	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

func TestSanitizeArg_Valid(t *testing.T) {
	got, err := SanitizeArg("  hello  ", "arg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestSanitizeArg_Empty(t *testing.T) {
	_, err := SanitizeArg("", "myarg")
	if err == nil {
		t.Fatal("expected error for empty arg")
	}
}

func TestSanitizeArg_Whitespace(t *testing.T) {
	_, err := SanitizeArg("   ", "myarg")
	if err == nil {
		t.Fatal("expected error for whitespace-only arg")
	}
}

func TestResolvePackage_FromEnv(t *testing.T) {
	t.Setenv(EnvPackage, "com.example.fromenv")

	cmd := &cobra.Command{}
	cmd.Flags().String("package", "", "")

	got, err := ResolvePackage(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "com.example.fromenv" {
		t.Errorf("expected %q, got %q", "com.example.fromenv", got)
	}
}

func TestGetOutputFormat_FromEnv(t *testing.T) {
	t.Setenv(EnvOutput, "yaml")

	cmd := &cobra.Command{}
	cmd.Flags().String("output", "", "")

	got := GetOutputFormat(cmd)
	if got != output.FormatYAML {
		t.Errorf("expected FormatYAML, got %v", got)
	}
}
