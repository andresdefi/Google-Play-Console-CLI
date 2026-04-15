package auth

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- MaskEmail ---

func TestMaskEmail_Long(t *testing.T) {
	email := "verylongusername@example.iam.gserviceaccount.com"
	masked := MaskEmail(email)
	if !strings.Contains(masked, "...") {
		t.Errorf("expected masked email to contain '...', got %q", masked)
	}
	if !strings.Contains(masked, "@example.iam.gserviceaccount.com") {
		t.Errorf("expected domain to be preserved, got %q", masked)
	}
	// Should be first 4 + "..." + last 4 of the local part.
	parts := strings.SplitN(masked, "@", 2)
	local := parts[0]
	if !strings.HasPrefix(local, "very") {
		t.Errorf("expected local part to start with 'very', got %q", local)
	}
	if !strings.HasSuffix(local, "name") {
		t.Errorf("expected local part to end with 'name', got %q", local)
	}
}

func TestMaskEmail_Short(t *testing.T) {
	// Local part <= 8 chars, should not be masked.
	email := "short@example.com"
	masked := MaskEmail(email)
	if masked != "short@example.com" {
		t.Errorf("expected 'short@example.com', got %q", masked)
	}
}

func TestMaskEmail_NoAt(t *testing.T) {
	// No @ sign, return as-is.
	input := "noemail"
	masked := MaskEmail(input)
	if masked != "noemail" {
		t.Errorf("expected 'noemail', got %q", masked)
	}
}

func TestMaskEmail_ExactlyEight(t *testing.T) {
	// Local part exactly 8 chars - should NOT be masked (len > 8 triggers masking).
	email := "abcdefgh@example.com"
	masked := MaskEmail(email)
	if masked != "abcdefgh@example.com" {
		t.Errorf("expected no masking for 8-char local, got %q", masked)
	}
}

func TestMaskEmail_NineChars(t *testing.T) {
	// Local part 9 chars - SHOULD be masked.
	email := "abcdefghi@example.com"
	masked := MaskEmail(email)
	if !strings.Contains(masked, "...") {
		t.Errorf("expected masking for 9-char local, got %q", masked)
	}
	expected := "abcd...fghi@example.com"
	if masked != expected {
		t.Errorf("expected %q, got %q", expected, masked)
	}
}

// --- Logout ---

func TestLogout_NoCredentials(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// Logout should not error even if nothing is stored.
	err := Logout()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// --- GetToken ---

func TestGetToken_NotAuthenticated(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	_, err := GetToken()
	if err == nil {
		t.Fatal("expected error when not authenticated")
	}
	if !strings.Contains(err.Error(), "not authenticated") {
		t.Errorf("expected 'not authenticated' error, got %q", err.Error())
	}
}

// --- GetStatus ---

func TestGetStatus_NotAuthenticated(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	_, err := GetStatus()
	if err == nil {
		t.Fatal("expected error when not authenticated")
	}
	if !strings.Contains(err.Error(), "not authenticated") {
		t.Errorf("expected 'not authenticated' error, got %q", err.Error())
	}
}

// --- Login ---

func TestLogin_InvalidKeyFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	_, err := Login("/nonexistent/key.json")
	if err == nil {
		t.Fatal("expected error for nonexistent key file")
	}
	if !strings.Contains(err.Error(), "could not read key file") {
		t.Errorf("expected 'could not read key file' error, got %q", err.Error())
	}
}

func TestLogin_InvalidJSONKeyFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	keyFile := filepath.Join(tmp, "bad-key.json")
	_ = os.WriteFile(keyFile, []byte("not valid json at all"), 0o600)

	_, err := Login(keyFile)
	if err == nil {
		t.Fatal("expected error for invalid key file")
	}
	if !strings.Contains(err.Error(), "invalid service account key file") {
		t.Errorf("expected 'invalid service account key file' error, got %q", err.Error())
	}
}

func TestGetToken_InvalidStoredCredentials(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// Write a config pointing to a key file with invalid content.
	dir := filepath.Join(tmp, ".gpc")
	_ = os.MkdirAll(dir, 0o700)
	_ = os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"key_file_path":"`+filepath.Join(tmp, "bad.json")+`"}`), 0o600)
	_ = os.WriteFile(filepath.Join(tmp, "bad.json"), []byte("not json"), 0o600)

	_, err := GetToken()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid stored credentials") {
		t.Errorf("expected 'invalid stored credentials' error, got %q", err.Error())
	}
}
