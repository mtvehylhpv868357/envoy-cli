package clone_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/clone"
	"github.com/user/envoy-cli/internal/profile"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "clone-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeStore(t *testing.T) *profile.Store {
	t.Helper()
	store, err := profile.LoadStore(filepath.Join(tempDir(t), "profiles"))
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func TestClone_Basic(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("prod", map[string]string{"DB": "prod-db", "PORT": "5432"})

	if err := clone.Profile(store, "prod", "staging", clone.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, err := store.Get("staging")
	if err != nil {
		t.Fatalf("cloned profile not found: %v", err)
	}
	if p.Vars["DB"] != "prod-db" {
		t.Errorf("expected DB=prod-db, got %q", p.Vars["DB"])
	}
}

func TestClone_WithOverrides(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("prod", map[string]string{"DB": "prod-db", "PORT": "5432"})

	opts := clone.DefaultOptions()
	opts.Overrides = map[string]string{"DB": "staging-db"}

	if err := clone.Profile(store, "prod", "staging", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, _ := store.Get("staging")
	if p.Vars["DB"] != "staging-db" {
		t.Errorf("expected DB=staging-db, got %q", p.Vars["DB"])
	}
	if p.Vars["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", p.Vars["PORT"])
	}
}

func TestClone_SourceNotFound(t *testing.T) {
	store := makeStore(t)
	err := clone.Profile(store, "ghost", "copy", clone.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing source, got nil")
	}
}

func TestClone_DestinationExists(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("a", map[string]string{"X": "1"})
	_ = store.Add("b", map[string]string{"X": "2"})

	err := clone.Profile(store, "a", "b", clone.DefaultOptions())
	if err == nil {
		t.Fatal("expected error when destination already exists")
	}
}

func TestClone_SameName(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("env", map[string]string{"K": "v"})

	err := clone.Profile(store, "env", "env", clone.DefaultOptions())
	if err == nil {
		t.Fatal("expected error when src == dst")
	}
}

func TestClone_SetActive(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("prod", map[string]string{"ENV": "production"})

	opts := clone.DefaultOptions()
	opts.SetActive = true

	if err := clone.Profile(store, "prod", "local", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	active, err := store.Active()
	if err != nil {
		t.Fatalf("could not retrieve active profile: %v", err)
	}
	if active.Name != "local" {
		t.Errorf("expected active=local, got %q", active.Name)
	}
}
