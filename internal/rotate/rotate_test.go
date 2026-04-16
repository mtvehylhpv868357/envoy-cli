package rotate_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/profile"
	"github.com/user/envoy-cli/internal/rotate"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "rotate-test-*")
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

func upperTransform(_, v string) (string, error) {
	return strings.ToUpper(v), nil
}

func TestRotate_AllKeys(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{"foo": "bar", "baz": "qux"})

	changed, err := rotate.Profile(store, "dev", upperTransform, rotate.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changed) != 2 {
		t.Fatalf("expected 2 changed keys, got %d", len(changed))
	}

	p, _ := store.Get("dev")
	if p.Vars["foo"] != "BAR" {
		t.Errorf("expected BAR, got %s", p.Vars["foo"])
	}
}

func TestRotate_SpecificKeys(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{"foo": "bar", "baz": "qux"})

	opts := rotate.DefaultOptions()
	opts.Keys = []string{"foo"}
	changed, err := rotate.Profile(store, "dev", upperTransform, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changed) != 1 {
		t.Fatalf("expected 1 changed key, got %d", len(changed))
	}

	p, _ := store.Get("dev")
	if p.Vars["baz"] != "qux" {
		t.Errorf("baz should be unchanged, got %s", p.Vars["baz"])
	}
}

func TestRotate_DryRun(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{"foo": "bar"})

	opts := rotate.DefaultOptions()
	opts.DryRun = true
	_, err := rotate.Profile(store, "dev", upperTransform, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, _ := store.Get("dev")
	if p.Vars["foo"] != "bar" {
		t.Errorf("dry-run should not persist changes, got %s", p.Vars["foo"])
	}
}

func TestRotate_ProfileNotFound(t *testing.T) {
	store := makeStore(t)
	_, err := rotate.Profile(store, "missing", upperTransform, rotate.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestRotate_EmptyName(t *testing.T) {
	store := makeStore(t)
	_, err := rotate.Profile(store, "", upperTransform, rotate.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}
