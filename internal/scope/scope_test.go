package scope_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envoy-cli/internal/scope"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "scope-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestNewStore_Empty(t *testing.T) {
	dir := tempDir(t)
	s, err := scope.NewStore(filepath.Join(dir, "scope.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.List()) != 0 {
		t.Errorf("expected empty store, got %d bindings", len(s.List()))
	}
}

func TestSet_And_Get(t *testing.T) {
	dir := tempDir(t)
	s, _ := scope.NewStore(filepath.Join(dir, "scope.json"))

	if err := s.Set("/projects/myapp", "production"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	if got := s.Get("/projects/myapp"); got != "production" {
		t.Errorf("expected production, got %q", got)
	}
}

func TestSet_Overwrites_Existing(t *testing.T) {
	dir := tempDir(t)
	s, _ := scope.NewStore(filepath.Join(dir, "scope.json"))

	_ = s.Set("/projects/myapp", "staging")
	_ = s.Set("/projects/myapp", "production")

	if got := s.Get("/projects/myapp"); got != "production" {
		t.Errorf("expected production after overwrite, got %q", got)
	}
	if len(s.List()) != 1 {
		t.Errorf("expected 1 binding, got %d", len(s.List()))
	}
}

func TestGet_NotFound(t *testing.T) {
	dir := tempDir(t)
	s, _ := scope.NewStore(filepath.Join(dir, "scope.json"))

	if got := s.Get("/nonexistent"); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestResolve_WalksUp(t *testing.T) {
	dir := tempDir(t)
	s, _ := scope.NewStore(filepath.Join(dir, "scope.json"))

	_ = s.Set("/projects/myapp", "development")

	got := s.Resolve("/projects/myapp/src/handlers")
	if got != "development" {
		t.Errorf("expected development from ancestor, got %q", got)
	}
}

func TestResolve_NoMatch(t *testing.T) {
	dir := tempDir(t)
	s, _ := scope.NewStore(filepath.Join(dir, "scope.json"))

	if got := s.Resolve("/unrelated/path"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestRemove(t *testing.T) {
	dir := tempDir(t)
	s, _ := scope.NewStore(filepath.Join(dir, "scope.json"))

	_ = s.Set("/projects/myapp", "production")
	if err := s.Remove("/projects/myapp"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	if got := s.Get("/projects/myapp"); got != "" {
		t.Errorf("expected empty after remove, got %q", got)
	}
}

func TestPersistence(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "scope.json")

	s1, _ := scope.NewStore(path)
	_ = s1.Set("/projects/alpha", "staging")

	s2, err := scope.NewStore(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if got := s2.Get("/projects/alpha"); got != "staging" {
		t.Errorf("expected staging after reload, got %q", got)
	}
}
