package snapshot_test

import (
	"os"
	"testing"

	"github.com/yourorg/envoy-cli/internal/snapshot"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "snapshot-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNewStore(t *testing.T) {
	dir := tempDir(t)
	_, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
}

func TestSaveAndLoad(t *testing.T) {
	store, _ := snapshot.NewStore(tempDir(t))
	snap := snapshot.Snapshot{
		Name:    "my-snap",
		Profile: "dev",
		Vars:    map[string]string{"FOO": "bar", "PORT": "8080"},
	}
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load("my-snap")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Profile != "dev" {
		t.Errorf("Profile: got %q, want %q", got.Profile, "dev")
	}
	if got.Vars["FOO"] != "bar" {
		t.Errorf("Vars[FOO]: got %q, want %q", got.Vars["FOO"], "bar")
	}
	if got.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestLoad_NotFound(t *testing.T) {
	store, _ := snapshot.NewStore(tempDir(t))
	_, err := store.Load("ghost")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestList(t *testing.T) {
	store, _ := snapshot.NewStore(tempDir(t))
	for _, name := range []string{"alpha", "beta", "gamma"} {
		_ = store.Save(snapshot.Snapshot{Name: name, Profile: "test", Vars: map[string]string{}})
	}
	names, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("List: got %d names, want 3", len(names))
	}
}

func TestDelete(t *testing.T) {
	store, _ := snapshot.NewStore(tempDir(t))
	_ = store.Save(snapshot.Snapshot{Name: "tmp", Profile: "p", Vars: map[string]string{}})
	if err := store.Delete("tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	names, _ := store.List()
	if len(names) != 0 {
		t.Errorf("expected 0 snapshots after delete, got %d", len(names))
	}
}

func TestDelete_NotFound(t *testing.T) {
	store, _ := snapshot.NewStore(tempDir(t))
	if err := store.Delete("nope"); err == nil {
		t.Fatal("expected error deleting non-existent snapshot")
	}
}
