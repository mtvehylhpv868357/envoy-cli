package prune_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/user/envoy-cli/internal/profile"
	"github.com/user/envoy-cli/internal/prune"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "prune-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeStore(t *testing.T) (*profile.Store, string) {
	t.Helper()
	dir := tempDir(t)
	path := filepath.Join(dir, "profiles")
	store, err := profile.LoadStore(path)
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	return store, path
}

func TestPrune_EmptyValues(t *testing.T) {
	store, _ := makeStore(t)
	_ = store.Add("dev", map[string]string{"FOO": "bar", "EMPTY": "", "BLANK": "   "})

	opts := prune.DefaultOptions()
	opts.EmptyValues = true

	res, err := prune.Profile(store, "dev", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d: %v", len(res.Removed), res.Removed)
	}
	if res.Kept != 1 {
		t.Errorf("expected 1 kept, got %d", res.Kept)
	}
}

func TestPrune_KeyPattern(t *testing.T) {
	store, _ := makeStore(t)
	_ = store.Add("prod", map[string]string{"DEBUG_A": "1", "DEBUG_B": "2", "APP_URL": "http://x"})

	opts := prune.DefaultOptions()
	opts.KeyPattern = "^DEBUG_"

	res, err := prune.Profile(store, "prod", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sort.Strings(res.Removed)
	if len(res.Removed) != 2 || res.Removed[0] != "DEBUG_A" || res.Removed[1] != "DEBUG_B" {
		t.Errorf("unexpected removed: %v", res.Removed)
	}
	if res.Kept != 1 {
		t.Errorf("expected 1 kept, got %d", res.Kept)
	}
}

func TestPrune_DryRun_DoesNotSave(t *testing.T) {
	store, _ := makeStore(t)
	_ = store.Add("staging", map[string]string{"KEY": "", "KEEP": "val"})

	opts := prune.DefaultOptions()
	opts.EmptyValues = true
	opts.DryRun = true

	res, err := prune.Profile(store, "staging", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 1 {
		t.Errorf("expected 1 in dry-run removed list, got %d", len(res.Removed))
	}

	// Reload and confirm original data unchanged.
	vars, err := store.Get("staging")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if _, ok := vars["KEY"]; !ok {
		t.Error("dry-run should not have removed KEY from store")
	}
}

func TestPrune_EmptyName_Errors(t *testing.T) {
	store, _ := makeStore(t)
	_, err := prune.Profile(store, "", prune.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty profile name")
	}
}

func TestPrune_InvalidPattern_Errors(t *testing.T) {
	store, _ := makeStore(t)
	_ = store.Add("x", map[string]string{"A": "1"})

	opts := prune.DefaultOptions()
	opts.KeyPattern = "[invalid"

	_, err := prune.Profile(store, "x", opts)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}
