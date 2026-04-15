package rename_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envoy-cli/internal/profile"
	"github.com/yourusername/envoy-cli/internal/rename"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "rename-test-*")
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

func TestRename_Basic(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{"FOO": "bar"})

	if err := rename.Profile(store, "dev", "staging", rename.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := store.Get("staging"); err != nil {
		t.Error("expected staging to exist")
	}
	if _, err := store.Get("dev"); err == nil {
		t.Error("expected dev to be deleted")
	}
}

func TestRename_SameName(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{})

	err := rename.Profile(store, "dev", "dev", rename.DefaultOptions())
	if err != rename.ErrSameName {
		t.Errorf("expected ErrSameName, got %v", err)
	}
}

func TestRename_SourceNotFound(t *testing.T) {
	store := makeStore(t)

	err := rename.Profile(store, "ghost", "real", rename.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestRename_DestinationExists_NoOverwrite(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{"A": "1"})
	_ = store.Add("prod", map[string]string{"A": "2"})

	err := rename.Profile(store, "dev", "prod", rename.DefaultOptions())
	if err == nil {
		t.Fatal("expected error when destination exists")
	}
}

func TestRename_DestinationExists_Overwrite(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{"A": "1"})
	_ = store.Add("prod", map[string]string{"A": "2"})

	opts := rename.Options{Overwrite: true}
	if err := rename.Profile(store, "dev", "prod", opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vars, err := store.Get("prod")
	if err != nil {
		t.Fatal("prod should exist")
	}
	if vars["A"] != "1" {
		t.Errorf("expected overwritten value A=1, got %s", vars["A"])
	}
}

func TestRename_PreservesActiveProfile(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{})
	_ = store.SetActive("dev")

	_ = rename.Profile(store, "dev", "staging", rename.DefaultOptions())

	active, err := store.ActiveName()
	if err != nil {
		t.Fatal("expected an active profile")
	}
	if active != "staging" {
		t.Errorf("expected active=staging, got %s", active)
	}
}
