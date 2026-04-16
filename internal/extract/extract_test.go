package extract_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy/internal/extract"
	"github.com/envoy-cli/envoy/internal/profile"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "extract-test-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeStore(t *testing.T) (*profile.Store, string) {
	t.Helper()
	dir := tempDir(t)
	path := filepath.Join(dir, "profiles.json")
	store, err := profile.LoadStore(path)
	if err != nil {
		t.Fatalf("load store: %v", err)
	}
	return store, path
}

func TestExtract_ByKeys(t *testing.T) {
	store, _ := makeStore(t)
	_ = store.Add("prod", map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "APP_ENV": "production"})

	opts := extract.DefaultOptions()
	opts.Keys = []string{"DB_HOST", "DB_PORT"}

	got, err := extract.Profile(store, "prod", "db-only", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(got))
	}
	if got["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST = %q, want %q", got["DB_HOST"], "localhost")
	}
}

func TestExtract_ByPrefix(t *testing.T) {
	store, _ := makeStore(t)
	_ = store.Add("prod", map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "APP_ENV": "production"})

	opts := extract.DefaultOptions()
	opts.Prefix = "DB_"

	got, err := extract.Profile(store, "prod", "db-vars", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(got))
	}
	if _, ok := got["APP_ENV"]; ok {
		t.Error("APP_ENV should not be in extracted vars")
	}
}

func TestExtract_DestinationExists_NoOverwrite(t *testing.T) {
	store, _ := makeStore(t)
	_ = store.Add("prod", map[string]string{"DB_HOST": "localhost"})
	_ = store.Add("db-only", map[string]string{"DB_HOST": "old"})

	opts := extract.DefaultOptions()
	opts.Keys = []string{"DB_HOST"}

	_, err := extract.Profile(store, "prod", "db-only", opts)
	if err == nil {
		t.Fatal("expected error when destination exists and Overwrite is false")
	}
}

func TestExtract_DestinationExists_WithOverwrite(t *testing.T) {
	store, _ := makeStore(t)
	_ = store.Add("prod", map[string]string{"DB_HOST": "newhost"})
	_ = store.Add("db-only", map[string]string{"DB_HOST": "old"})

	opts := extract.DefaultOptions()
	opts.Keys = []string{"DB_HOST"}
	opts.Overwrite = true

	got, err := extract.Profile(store, "prod", "db-only", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_HOST"] != "newhost" {
		t.Errorf("DB_HOST = %q, want %q", got["DB_HOST"], "newhost")
	}
}

func TestExtract_MissingKey_ReturnsError(t *testing.T) {
	store, _ := makeStore(t)
	_ = store.Add("prod", map[string]string{"APP_ENV": "production"})

	opts := extract.DefaultOptions()
	opts.Keys = []string{"MISSING_KEY"}

	_, err := extract.Profile(store, "prod", "subset", opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestExtract_SameName_ReturnsError(t *testing.T) {
	store, _ := makeStore(t)
	_ = store.Add("prod", map[string]string{"DB_HOST": "localhost"})

	opts := extract.DefaultOptions()
	opts.Keys = []string{"DB_HOST"}

	_, err := extract.Profile(store, "prod", "prod", opts)
	if err == nil {
		t.Fatal("expected error when src and dst are the same")
	}
}
