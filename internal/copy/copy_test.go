package copy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envoy-cli/internal/copy"
	"github.com/yourusername/envoy-cli/internal/profile"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "copy-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeStore(t *testing.T) *profile.Store {
	t.Helper()
	store, err := profile.LoadStore(filepath.Join(tempDir(t), "profiles"))
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	return store
}

// mustAdd adds a profile to the store and fails the test on error.
func mustAdd(t *testing.T, store *profile.Store, name string, vars map[string]string) {
	t.Helper()
	if err := store.Add(name, vars); err != nil {
		t.Fatalf("Add %q: %v", name, err)
	}
}

func TestCopy_Basic(t *testing.T) {
	store := makeStore(t)
	mustAdd(t, store, "dev", map[string]string{"FOO": "bar", "PORT": "8080"})

	if err := copy.Profile(store, "dev", "staging", copy.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, err := store.Get("staging")
	if err != nil {
		t.Fatalf("Get staging: %v", err)
	}
	if p.Vars["FOO"] != "bar" || p.Vars["PORT"] != "8080" {
		t.Errorf("vars mismatch: %v", p.Vars)
	}
}

func TestCopy_SameName(t *testing.T) {
	store := makeStore(t)
	mustAdd(t, store, "dev", map[string]string{"X": "1"})

	err := copy.Profile(store, "dev", "dev", copy.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for same src/dst names")
	}
}

func TestCopy_SourceNotFound(t *testing.T) {
	store := makeStore(t)

	err := copy.Profile(store, "nonexistent", "dst", copy.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestCopy_DestinationExists_NoOverwrite(t *testing.T) {
	store := makeStore(t)
	mustAdd(t, store, "dev", map[string]string{"A": "1"})
	mustAdd(t, store, "prod", map[string]string{"A": "2"})

	err := copy.Profile(store, "dev", "prod", copy.DefaultOptions())
	if err == nil {
		t.Fatal("expected error when destination exists and overwrite is false")
	}
}

func TestCopy_DestinationExists_WithOverwrite(t *testing.T) {
	store := makeStore(t)
	mustAdd(t, store, "dev", map[string]string{"A": "1"})
	mustAdd(t, store, "prod", map[string]string{"A": "2"})

	opts := copy.Options{Overwrite: true}
	if err := copy.Profile(store, "dev", "prod", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, _ := store.Get("prod")
	if p.Vars["A"] != "1" {
		t.Errorf("expected A=1, got %q", p.Vars["A"])
	}
}

func TestCopy_IsolatesVars(t *testing.T) {
	store := makeStore(t)
	mustAdd(t, store, "src", map[string]string{"K": "original"})

	_ = copy.Profile(store, "src", "dst", copy.DefaultOptions())

	// Mutate source after copy — destination should be unaffected.
	srcP, _ := store.Get("src")
	srcP.Vars["K"] = "mutated"

	dstP, _ := store.Get("dst")
	if dstP.Vars["K"] != "original" {
		t.Errorf("copy was not isolated: got %q", dstP.Vars["K"])
	}
}
