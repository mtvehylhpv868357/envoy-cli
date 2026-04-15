package rollback_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/envoy-cli/internal/history"
	"github.com/yourusername/envoy-cli/internal/profile"
	"github.com/yourusername/envoy-cli/internal/rollback"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "rollback-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestRollback_NoHistory(t *testing.T) {
	dir := tempDir(t)
	store, _ := profile.LoadStore(filepath.Join(dir, "profiles"))
	hist, _ := history.NewStore(filepath.Join(dir, "history"))

	_, err := rollback.Profile("dev", store, hist, rollback.DefaultOptions())
	if !errors.Is(err, rollback.ErrNoHistory) {
		t.Fatalf("expected ErrNoHistory, got %v", err)
	}
}

func TestRollback_IndexOutOfRange(t *testing.T) {
	dir := tempDir(t)
	store, _ := profile.LoadStore(filepath.Join(dir, "profiles"))
	hist, _ := history.NewStore(filepath.Join(dir, "history"))

	_ = hist.Record("dev", map[string]string{"A": "1"}, time.Now())

	_, err := rollback.Profile("dev", store, hist, rollback.Options{Index: 5})
	if !errors.Is(err, rollback.ErrIndexOutOfRange) {
		t.Fatalf("expected ErrIndexOutOfRange, got %v", err)
	}
}

func TestRollback_RestoresMostRecent(t *testing.T) {
	dir := tempDir(t)
	store, _ := profile.LoadStore(filepath.Join(dir, "profiles"))
	hist, _ := history.NewStore(filepath.Join(dir, "history"))

	old := map[string]string{"KEY": "old"}
	new_ := map[string]string{"KEY": "new", "EXTRA": "yes"}

	_ = hist.Record("dev", old, time.Now().Add(-time.Minute))
	_ = hist.Record("dev", new_, time.Now())

	res, err := rollback.Profile("dev", store, hist, rollback.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.VarsCount != 2 {
		t.Errorf("expected 2 vars, got %d", res.VarsCount)
	}

	vars, err := store.Get("dev")
	if err != nil {
		t.Fatalf("store.Get: %v", err)
	}
	if vars["KEY"] != "new" {
		t.Errorf("expected KEY=new, got %q", vars["KEY"])
	}
}

func TestRollback_EmptyName(t *testing.T) {
	dir := tempDir(t)
	store, _ := profile.LoadStore(filepath.Join(dir, "profiles"))
	hist, _ := history.NewStore(filepath.Join(dir, "history"))

	_, err := rollback.Profile("", store, hist, rollback.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty profile name")
	}
}
