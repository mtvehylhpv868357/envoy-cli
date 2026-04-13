package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envoy-cli/internal/config"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envoy-config-*")
	if err != nil {
		t.Fatalf("failed to create tempn	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig(Output != true {
		t.Errorf("expected ColorOutput=true, got %v", cfg.ColorOutput)
	}
	if cfg.AutoExport != false {
		t.Errorf("expected AutoExport=false, got %v", cfg.AutoExport)
	}
}

func TestLoad_NoFile_ReturnsDefaults(t *testing.T) {
	dir := tempDir(t)
	cfg, err := config.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.ColorOutput != true {
		t.Errorf("expected default ColorOutput=true")
	}
}

func TestSRoundTrip(t *testing.T) {
	dir := tempDir(t)

	orig := &config.Config{
		DefaultShell:",
		ProfilesDir:  "/custom/profiles",
		AutoExport:   true,
		ColorOutput:  false,
	}

	if err := config.Save(dir, orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := config.Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.DefaultShell != orig.DefaultShell {
		t.Errorf("DefaultShell: got %q, want %q", loaded.DefaultShell, orig.DefaultShell)
	}
	if loaded.ProfilesDir != orig.ProfilesDir {
		t.Errorf("ProfilesDir: got %q, want %q", loaded.ProfilesDir, orig.ProfilesDir)
	}
	if loaded.AutoExport != orig.AutoExport {
		t.Errorf("AutoExport: got %v, want %v", loaded.AutoExport, orig.AutoExport)
	}
	if loaded.ColorOutput != orig.ColorOutput {
		t.Errorf("ColorOutput: got %v, want %v", loaded.ColorOutput, orig.ColorOutput)
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	base := tempDir(t)
	nestedDir := filepath.Join(base, "deep", "nested")

	cfg := config.DefaultConfig()
	if err := config.Save(nestedDir, cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(nestedDir, "config.json")); os.IsNotExist(err) {
		t.Error("expected config.json to be created")
	}
}
