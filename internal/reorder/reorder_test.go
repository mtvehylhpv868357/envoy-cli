package reorder_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/internal/profile"
	"github.com/envoy-cli/internal/reorder"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "reorder-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func makeStore(t *testing.T) *profile.Store {
	t.Helper()
	st, err := profile.LoadStore(filepath.Join(tempDir(t), "profiles"))
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestKeys_Alpha(t *testing.T) {
	vars := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys := reorder.Keys(vars, reorder.Options{Strategy: reorder.StrategyAlpha})
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != want[i] {
			t.Fatalf("index %d: got %s want %s", i, k, want[i])
		}
	}
}

func TestKeys_AlphaDesc(t *testing.T) {
	vars := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys := reorder.Keys(vars, reorder.Options{Strategy: reorder.StrategyAlphaDesc})
	if keys[0] != "ZEBRA" {
		t.Fatalf("expected ZEBRA first, got %s", keys[0])
	}
}

func TestKeys_Custom(t *testing.T) {
	vars := map[string]string{"C": "3", "A": "1", "B": "2", "D": "4"}
	opts := reorder.Options{
		Strategy: reorder.StrategyCustom,
		Order:    []string{"B", "A"},
	}
	keys := reorder.Keys(vars, opts)
	if keys[0] != "B" || keys[1] != "A" {
		t.Fatalf("unexpected order: %v", keys)
	}
	// C and D should follow, sorted
	if keys[2] != "C" || keys[3] != "D" {
		t.Fatalf("unexpected tail order: %v", keys)
	}
}

func TestProfile_Reorder(t *testing.T) {
	st := makeStore(t)
	if err := st.Add("dev", map[string]string{"Z": "z", "A": "a"}); err != nil {
		t.Fatal(err)
	}
	_, err := reorder.Profile(st, "dev", reorder.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
}

func TestProfile_EmptyName(t *testing.T) {
	st := makeStore(t)
	_, err := reorder.Profile(st, "", reorder.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestProfile_NotFound(t *testing.T) {
	st := makeStore(t)
	_, err := reorder.Profile(st, "ghost", reorder.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}
