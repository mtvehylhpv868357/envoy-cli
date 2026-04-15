package merge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/internal/merge"
	"github.com/envoy-cli/internal/profile"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "merge-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
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

func TestMerge_NoConflicts(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("base", map[string]string{"A": "1", "B": "2"})
	_ = store.Add("extra", map[string]string{"C": "3"})

	res, err := merge.Profiles(store, "base", "extra", merge.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
	if res.Merged["C"] != "3" {
		t.Errorf("expected C=3, got %q", res.Merged["C"])
	}
}

func TestMerge_StrategyTheirs(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("base", map[string]string{"KEY": "old"})
	_ = store.Add("src", map[string]string{"KEY": "new"})

	opts := merge.DefaultOptions() // StrategyTheirs
	res, err := merge.Profiles(store, "base", "src", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "new" {
		t.Errorf("expected KEY=new, got %q", res.Merged["KEY"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "KEY" {
		t.Errorf("expected conflict on KEY, got %v", res.Conflicts)
	}
}

func TestMerge_StrategyOurs(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("base", map[string]string{"KEY": "old"})
	_ = store.Add("src", map[string]string{"KEY": "new"})

	opts := merge.Options{Strategy: merge.StrategyOurs}
	res, err := merge.Profiles(store, "base", "src", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "old" {
		t.Errorf("expected KEY=old, got %q", res.Merged["KEY"])
	}
}

func TestMerge_SourceNotFound(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("base", map[string]string{"A": "1"})

	_, err := merge.Profiles(store, "base", "missing", merge.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing source profile")
	}
}

func TestMerge_Overwrite(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("base", map[string]string{"A": "1"})
	_ = store.Add("extra", map[string]string{"B": "2"})

	opts := merge.Options{Strategy: merge.StrategyTheirs, Overwrite: true}
	_, err := merge.Profiles(store, "base", "extra", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := store.Get("base")
	if err != nil {
		t.Fatalf("Get base: %v", err)
	}
	if updated.Vars["B"] != "2" {
		t.Errorf("expected B=2 in base after overwrite, got %q", updated.Vars["B"])
	}
}
