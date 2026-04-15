package trim_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/profile"
	"github.com/envoy-cli/envoy-cli/internal/trim"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "trim-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeStore(t *testing.T, vars map[string]string) (*profile.Store, string) {
	t.Helper()
	dir := tempDir(t)
	store, err := profile.LoadStore(filepath.Join(dir, "profiles"))
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	p := &profile.Profile{Name: "test", Vars: vars}
	if err := store.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return store, dir
}

func TestTrim_ExplicitKeys(t *testing.T) {
	store, _ := makeStore(t, map[string]string{"A": "1", "B": "2", "C": "3"})
	opts := trim.DefaultOptions()
	opts.Keys = []string{"A", "C"}

	res, err := trim.Profile(store, "test", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sort.Strings(res.Removed)
	if len(res.Removed) != 2 || res.Removed[0] != "A" || res.Removed[1] != "C" {
		t.Errorf("removed = %v, want [A C]", res.Removed)
	}

	p, _ := store.Get("test")
	if _, ok := p.Vars["B"]; !ok {
		t.Error("key B should remain")
	}
	if _, ok := p.Vars["A"]; ok {
		t.Error("key A should have been removed")
	}
}

func TestTrim_Prefix(t *testing.T) {
	store, _ := makeStore(t, map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "DB_URL": "postgres://"})
	opts := trim.DefaultOptions()
	opts.Prefix = "APP_"

	res, err := trim.Profile(store, "test", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	p, _ := store.Get("test")
	if _, ok := p.Vars["DB_URL"]; !ok {
		t.Error("DB_URL should remain")
	}
}

func TestTrim_EmptyValues(t *testing.T) {
	store, _ := makeStore(t, map[string]string{"KEEP": "value", "GONE": "", "ALSO_GONE": ""})
	opts := trim.DefaultOptions()
	opts.EmptyValues = true

	res, err := trim.Profile(store, "test", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
}

func TestTrim_DryRun_DoesNotSave(t *testing.T) {
	store, _ := makeStore(t, map[string]string{"X": "1", "Y": "2"})
	opts := trim.DefaultOptions()
	opts.DryRun = true
	opts.Keys = []string{"X"}

	res, err := trim.Profile(store, "test", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 {
		t.Errorf("dry-run removed count = %d, want 1", len(res.Removed))
	}
	p, _ := store.Get("test")
	if _, ok := p.Vars["X"]; !ok {
		t.Error("dry-run should not persist removal of X")
	}
}

func TestTrim_EmptyName_ReturnsError(t *testing.T) {
	store, _ := makeStore(t, map[string]string{})
	_, err := trim.Profile(store, "", trim.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty name")
	}
}
