package squash_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/profile"
	"github.com/envoy-cli/envoy-cli/internal/squash"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "squash-test-*")
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

func TestSquash_BasicMerge(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("base", map[string]string{"A": "1", "B": "2"})
	_ = store.Add("override", map[string]string{"B": "99", "C": "3"})

	result, err := squash.Profiles(store, []string{"base", "override"}, squash.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["A"] != "1" || result["B"] != "99" || result["C"] != "3" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestSquash_NoOverwrite(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("first", map[string]string{"X": "first"})
	_ = store.Add("second", map[string]string{"X": "second"})

	opts := squash.Options{Overwrite: false}
	result, err := squash.Profiles(store, []string{"first", "second"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["X"] != "first" {
		t.Errorf("expected first value to win, got %q", result["X"])
	}
}

func TestSquash_ProfileNotFound(t *testing.T) {
	store := makeStore(t)
	_, err := squash.Profiles(store, []string{"missing"}, squash.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestSquash_EmptyNames(t *testing.T) {
	store := makeStore(t)
	_, err := squash.Profiles(store, []string{}, squash.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty names slice")
	}
}
