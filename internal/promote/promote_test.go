package promote_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envoy-cli/internal/profile"
	"github.com/yourusername/envoy-cli/internal/promote"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "promote-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeStore(t *testing.T, dir string) *profile.Store {
	t.Helper()
	store, err := profile.LoadStore(filepath.Join(dir, "profiles"))
	if err != nil {
		t.Fatalf("makeStore: %v", err)
	}
	return store
}

func TestPromote_Basic(t *testing.T) {
	dir := tempDir(t)
	store := makeStore(t, dir)

	_ = store.Add("staging", map[string]string{"APP_ENV": "staging", "PORT": "8080"})

	opts := promote.DefaultOptions()
	opts.StorePath = filepath.Join(dir, "profiles")

	res, err := promote.Profile("staging", "production", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Source != "staging" || res.Destination != "production" {
		t.Errorf("unexpected result: %+v", res)
	}

	// Reload store and verify vars were copied.
	store2, _ := profile.LoadStore(filepath.Join(dir, "profiles"))
	vars, err := store2.Get("production")
	if err != nil {
		t.Fatalf("production profile not found: %v", err)
	}
	if vars["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", vars["APP_ENV"])
	}
}

func TestPromote_OverwriteFalse_Errors(t *testing.T) {
	dir := tempDir(t)
	store := makeStore(t, dir)
	_ = store.Add("staging", map[string]string{"X": "1"})
	_ = store.Add("production", map[string]string{"X": "2"})

	opts := promote.DefaultOptions()
	opts.StorePath = filepath.Join(dir, "profiles")

	_, err := promote.Profile("staging", "production", opts)
	if err == nil {
		t.Fatal("expected error when destination exists and overwrite=false")
	}
}

func TestPromote_OverwriteTrue(t *testing.T) {
	dir := tempDir(t)
	store := makeStore(t, dir)
	_ = store.Add("staging", map[string]string{"X": "new"})
	_ = store.Add("production", map[string]string{"X": "old"})

	opts := promote.DefaultOptions()
	opts.StorePath = filepath.Join(dir, "profiles")
	opts.Overwrite = true

	_, err := promote.Profile("staging", "production", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	store2, _ := profile.LoadStore(filepath.Join(dir, "profiles"))
	vars, _ := store2.Get("production")
	if vars["X"] != "new" {
		t.Errorf("expected X=new after overwrite, got %q", vars["X"])
	}
}

func TestPromote_Activate(t *testing.T) {
	dir := tempDir(t)
	store := makeStore(t, dir)
	_ = store.Add("staging", map[string]string{"ENV": "stg"})

	opts := promote.DefaultOptions()
	opts.StorePath = filepath.Join(dir, "profiles")
	opts.Activate = true

	res, err := promote.Profile("staging", "production", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Activated {
		t.Error("expected Activated=true")
	}
}

func TestPromote_SameName_Errors(t *testing.T) {
	dir := tempDir(t)
	store := makeStore(t, dir)
	_ = store.Add("staging", map[string]string{"X": "1"})

	opts := promote.DefaultOptions()
	opts.StorePath = filepath.Join(dir, "profiles")

	_, err := promote.Profile("staging", "staging", opts)
	if err == nil {
		t.Fatal("expected error when source and destination are the same")
	}
}
