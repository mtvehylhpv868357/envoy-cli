package completion_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/completion"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "completion-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

func TestProfileNames_Empty(t *testing.T) {
	dir := tempDir(t)
	names := completion.ProfileNames(dir)
	if len(names) != 0 {
		t.Errorf("expected 0 names, got %d", len(names))
	}
}

func TestProfileNames_ReturnsJsonFiles(t *testing.T) {
	dir := tempDir(t)
	writeFile(t, dir, "dev.json", "{}")
	writeFile(t, dir, "prod.json", "{}")
	writeFile(t, dir, "notes.txt", "ignore me")

	names := completion.ProfileNames(dir)
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d: %v", len(names), names)
	}
	for _, n := range names {
		if n != "dev" && n != "prod" {
			t.Errorf("unexpected name: %s", n)
		}
	}
}

func TestProfileNames_NoDirectory(t *testing.T) {
	names := completion.ProfileNames("/nonexistent/path/xyz")
	if len(names) != 0 {
		t.Errorf("expected empty slice for missing dir, got %v", names)
	}
}

func TestEnvFiles_FindsDotEnv(t *testing.T) {
	dir := tempDir(t)
	writeFile(t, dir, ".env", "KEY=val")
	writeFile(t, dir, ".env.local", "KEY=local")
	writeFile(t, dir, "main.go", "package main")

	files := completion.EnvFiles(dir)
	if len(files) != 2 {
		t.Fatalf("expected 2 env files, got %d: %v", len(files), files)
	}
}

func TestFilterPrefix(t *testing.T) {
	items := []string{"dev", "dev-local", "prod", "staging"}

	got := completion.FilterPrefix(items, "dev")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(got), got)
	}

	all := completion.FilterPrefix(items, "")
	if len(all) != 4 {
		t.Errorf("empty prefix should return all, got %d", len(all))
	}

	none := completion.FilterPrefix(items, "zzz")
	if len(none) != 0 {
		t.Errorf("expected 0, got %d", len(none))
	}
}
