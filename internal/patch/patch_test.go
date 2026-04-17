package patch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/internal/patch"
	"github.com/envoy-cli/internal/profile"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "patch-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func makeStore(t *testing.T) *profile.Store {
	t.Helper()
	s, err := profile.LoadStore(filepath.Join(tempDir(t), "profiles"))
	if err != nil {
		t.Fatalf("makeStore: %v", err)
	}
	return s
}

func TestPatch_Upsert(t *testing.T) {
	s := makeStore(t)
	_ = s.Add("dev", map[string]string{"A": "1", "B": "2"})

	opts := patch.DefaultOptions()
	opts.Upsert = map[string]string{"B": "99", "C": "3"}

	result, err := patch.Profile(s, "dev", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["B"] != "99" {
		t.Errorf("expected B=99, got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
	if result["A"] != "1" {
		t.Errorf("expected A=1 unchanged, got %s", result["A"])
	}
}

func TestPatch_Delete(t *testing.T) {
	s := makeStore(t)
	_ = s.Add("dev", map[string]string{"X": "1", "Y": "2"})

	opts := patch.DefaultOptions()
	opts.Delete = []string{"X"}

	result, err := patch.Profile(s, "dev", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["X"]; ok {
		t.Error("expected X to be deleted")
	}
	if result["Y"] != "2" {
		t.Errorf("expected Y=2, got %s", result["Y"])
	}
}

func TestPatch_EmptyName(t *testing.T) {
	s := makeStore(t)
	_, err := patch.Profile(s, "", patch.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestPatch_ProfileNotFound(t *testing.T) {
	s := makeStore(t)
	_, err := patch.Profile(s, "ghost", patch.DefaultOptions())
	if err == nil {
		t.Error("expected error for missing profile")
	}
}
