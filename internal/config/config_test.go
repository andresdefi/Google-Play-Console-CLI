package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir, err := Dir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(tmp, ".gpc")
	if dir != expected {
		t.Errorf("expected %q, got %q", expected, dir)
	}
}

func TestPath(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	p, err := Path()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(tmp, ".gpc", "config.json")
	if p != expected {
		t.Errorf("expected %q, got %q", expected, p)
	}
}

func TestLoad_NoFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.KeyFilePath != "" {
		t.Errorf("expected empty KeyFilePath, got %q", cfg.KeyFilePath)
	}
	if cfg.PackageName != "" {
		t.Errorf("expected empty PackageName, got %q", cfg.PackageName)
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".gpc")
	_ = os.MkdirAll(dir, 0o700)

	data := `{"key_file_path":"/path/to/key.json","package_name":"com.example"}`
	_ = os.WriteFile(filepath.Join(dir, "config.json"), []byte(data), 0o600)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.KeyFilePath != "/path/to/key.json" {
		t.Errorf("expected KeyFilePath '/path/to/key.json', got %q", cfg.KeyFilePath)
	}
	if cfg.PackageName != "com.example" {
		t.Errorf("expected PackageName 'com.example', got %q", cfg.PackageName)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".gpc")
	_ = os.MkdirAll(dir, 0o700)
	_ = os.WriteFile(filepath.Join(dir, "config.json"), []byte("not json!"), 0o600)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSave_NewConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg := &Config{KeyFilePath: "/key.json", PackageName: "com.test"}
	if err := Save(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify file exists.
	p := filepath.Join(tmp, ".gpc", "config.json")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("could not read saved config: %v", err)
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("could not unmarshal saved config: %v", err)
	}
	if loaded.KeyFilePath != "/key.json" {
		t.Errorf("expected '/key.json', got %q", loaded.KeyFilePath)
	}
	if loaded.PackageName != "com.test" {
		t.Errorf("expected 'com.test', got %q", loaded.PackageName)
	}
}

func TestSave_OverwriteConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg1 := &Config{KeyFilePath: "/old.json"}
	_ = Save(cfg1)

	cfg2 := &Config{KeyFilePath: "/new.json"}
	if err := Save(cfg2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.KeyFilePath != "/new.json" {
		t.Errorf("expected '/new.json', got %q", loaded.KeyFilePath)
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg := &Config{PackageName: "com.test"}
	if err := Save(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dir := filepath.Join(tmp, ".gpc")
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory, got file")
	}
}

func TestClear_Exists(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	_ = Save(&Config{KeyFilePath: "/key.json"})

	if err := Clear(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p := filepath.Join(tmp, ".gpc", "config.json")
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("expected config file to be removed")
	}
}

func TestClear_NotExists(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// Should not error even if file doesn't exist.
	if err := Clear(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfig_RoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	original := &Config{
		KeyFilePath: "/path/to/sa.json",
		PackageName: "com.example.app",
	}
	if err := Save(original); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.KeyFilePath != original.KeyFilePath {
		t.Errorf("KeyFilePath mismatch: %q vs %q", loaded.KeyFilePath, original.KeyFilePath)
	}
	if loaded.PackageName != original.PackageName {
		t.Errorf("PackageName mismatch: %q vs %q", loaded.PackageName, original.PackageName)
	}
}

func TestConfig_Permissions(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	_ = Save(&Config{KeyFilePath: "/key.json"})

	// Check directory permissions.
	dir := filepath.Join(tmp, ".gpc")
	dirInfo, err := os.Stat(dir)
	if err != nil {
		t.Fatal(err)
	}
	dirPerm := dirInfo.Mode().Perm()
	if dirPerm != 0o700 {
		t.Errorf("expected dir perm 0700, got %04o", dirPerm)
	}

	// Check file permissions - the file is written via temp+rename.
	// The temp file is written with 0600.
	p := filepath.Join(dir, "config.json")
	fileInfo, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	filePerm := fileInfo.Mode().Perm()
	if filePerm != 0o600 {
		t.Errorf("expected file perm 0600, got %04o", filePerm)
	}
}

func TestConfig_AtomicWrite(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	_ = Save(&Config{KeyFilePath: "/key.json"})

	// After save, there should be no .tmp file remaining.
	p := filepath.Join(tmp, ".gpc", "config.json.tmp")
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("expected .tmp file to be cleaned up after atomic write")
	}
}

func TestConfig_EmptyFields(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg := &Config{}
	_ = Save(cfg)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.KeyFilePath != "" {
		t.Errorf("expected empty KeyFilePath, got %q", loaded.KeyFilePath)
	}
	if loaded.PackageName != "" {
		t.Errorf("expected empty PackageName, got %q", loaded.PackageName)
	}
}

func TestConfig_AllFields(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg := &Config{
		KeyFilePath: "/very/long/path/to/service-account.json",
		PackageName: "com.very.long.package.name.app",
	}
	_ = Save(cfg)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.KeyFilePath != cfg.KeyFilePath {
		t.Errorf("KeyFilePath: expected %q, got %q", cfg.KeyFilePath, loaded.KeyFilePath)
	}
	if loaded.PackageName != cfg.PackageName {
		t.Errorf("PackageName: expected %q, got %q", cfg.PackageName, loaded.PackageName)
	}
}
