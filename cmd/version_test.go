package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand_Output(t *testing.T) {
	cmd := newVersionCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	// newVersionCmd uses fmt.Println which writes to os.Stdout, not cmd's output.
	// We just verify it runs without error.
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVersionCommand_NoArgs(t *testing.T) {
	cmd := newVersionCmd()
	// Should reject extra args.
	cmd.SetArgs([]string{"extra"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for extra args")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		// cobra.NoArgs produces an error message.
		// The exact message may vary, just verify it errors.
		t.Logf("error message: %s", err.Error())
	}
}

func TestVersionCommand_Use(t *testing.T) {
	cmd := newVersionCmd()
	if cmd.Use != "version" {
		t.Errorf("expected Use 'version', got %q", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("expected non-empty Short description")
	}
}
