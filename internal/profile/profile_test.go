package profile

import (
	"os"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envoy-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestLoadStore_NewStore(t *testing.T) {
	store, err := LoadStore(tempDir(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Profiles) != 0 {
		t.Errorf("expected empty profiles, got %d", len(store.Profiles))
	}
}

func TestAddAndGet(t *testing.T) {
	store, _ := LoadStore(tempDir(t))
	p := &Profile{Name: "dev", Vars: map[string]string{"DB_URL": "localhost"}}
	store.Add(p)

	got, ok := store.Get("dev")
	if !ok {
		t.Fatal("expected profile to exist")
	}
	if got.Vars["DB_URL"] != "localhost" {
		t.Errorf("unexpected DB_URL: %s", got.Vars["DB_URL"])
	}
}

func TestDelete(t *testing.T) {
	store, _ := LoadStore(tempDir(t))
	store.Add(&Profile{Name: "staging", Vars: map[string]string{}})
	store.Active = "staging"

	if !store.Delete("staging") {
		t.Fatal("expected delete to return true")
	}
	if _, ok := store.Get("staging"); ok {
		t.Error("profile should have been deleted")
	}
	if store.Active != "" {
		t.Errorf("active should be cleared, got %q", store.Active)
	}
	if store.Delete("nonexistent") {
		t.Error("expected false for missing profile")
	}
}

func TestSetActive(t *testing.T) {
	store, _ := LoadStore(tempDir(t))
	store.Add(&Profile{Name: "prod", Vars: map[string]string{}})

	if err := store.SetActive("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.Active != "prod" {
		t.Errorf("expected active=prod, got %q", store.Active)
	}
	if err := store.SetActive("ghost"); err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestSaveAndReload(t *testing.T) {
	dir := tempDir(t)
	store, _ := LoadStore(dir)
	store.Add(&Profile{Name: "ci", Vars: map[string]string{"CI": "true"}})
	store.Active = "ci"

	if err := store.Save(); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	reloaded, err := LoadStore(dir)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if reloaded.Active != "ci" {
		t.Errorf("expected active=ci, got %q", reloaded.Active)
	}
	if p, ok := reloaded.Get("ci"); !ok || p.Vars["CI"] != "true" {
		t.Error("reloaded profile data mismatch")
	}
}
