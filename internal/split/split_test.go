package split_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/internal/profile"
	"github.com/envoy-cli/internal/split"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "split-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func makeStore(t *testing.T) *profile.Store {
	t.Helper()
	store, err := profile.LoadStore(filepath.Join(tempDir(t), "profiles"))
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func TestByPrefix_Basic(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("mono", map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_PORT": "8080",
	})

	opts := split.DefaultOptions()
	created, err := split.ByPrefix(store, "mono", []string{"DB_", "APP_"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(created) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(created))
	}

	db, err := store.Get("db")
	if err != nil {
		t.Fatal(err)
	}
	if db["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", db["HOST"])
	}
}

func TestByPrefix_StripPrefixFalse(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("mono", map[string]string{"DB_HOST": "pg"})

	opts := split.DefaultOptions()
	opts.StripPrefix = false
	_, err := split.ByPrefix(store, "mono", []string{"DB_"}, opts)
	if err != nil {
		t.Fatal(err)
	}
	db, _ := store.Get("db")
	if _, ok := db["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key when StripPrefix=false")
	}
}

func TestByPrefix_NoOverwrite(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("mono", map[string]string{"DB_HOST": "pg"})
	_ = store.Add("db", map[string]string{"HOST": "existing"})

	_, err := split.ByPrefix(store, "mono", []string{"DB_"}, split.DefaultOptions())
	if err == nil {
		t.Error("expected error when destination exists and Overwrite=false")
	}
}

func TestByPrefix_SourceNotFound(t *testing.T) {
	store := makeStore(t)
	_, err := split.ByPrefix(store, "missing", []string{"DB_"}, split.DefaultOptions())
	if err == nil {
		t.Error("expected error for missing source profile")
	}
}

func TestByPrefix_EmptyPrefixes(t *testing.T) {
	store := makeStore(t)
	_, err := split.ByPrefix(store, "mono", nil, split.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty prefixes")
	}
}
