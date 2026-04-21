package inherit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy/internal/inherit"
	"github.com/envoy-cli/envoy/internal/profile"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "inherit-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func makeStore(t *testing.T) *profile.Store {
	t.Helper()
	st, err := profile.LoadStore(filepath.Join(tempDir(t), "profiles"))
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	return st
}

func TestInherit_Basic(t *testing.T) {
	st := makeStore(t)
	_ = st.Add("base", map[string]string{"A": "1", "B": "2"})

	if err := inherit.Profile(st, "base", "child", nil, inherit.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vars, _ := st.Get("child")
	if vars["A"] != "1" || vars["B"] != "2" {
		t.Errorf("expected inherited vars, got %v", vars)
	}
}

func TestInherit_OverridesWin(t *testing.T) {
	st := makeStore(t)
	_ = st.Add("base", map[string]string{"A": "1", "B": "2"})

	overrides := map[string]string{"B": "override", "C": "3"}
	_ = inherit.Profile(st, "base", "child", overrides, inherit.DefaultOptions())

	vars, _ := st.Get("child")
	if vars["B"] != "override" {
		t.Errorf("expected override value for B, got %q", vars["B"])
	}
	if vars["C"] != "3" {
		t.Errorf("expected C=3, got %q", vars["C"])
	}
	if vars["A"] != "1" {
		t.Errorf("expected A=1 inherited, got %q", vars["A"])
	}
}

func TestInherit_DestinationExists_NoOverwrite(t *testing.T) {
	st := makeStore(t)
	_ = st.Add("base", map[string]string{"A": "1"})
	_ = st.Add("child", map[string]string{"X": "9"})

	err := inherit.Profile(st, "base", "child", nil, inherit.DefaultOptions())
	if err == nil {
		t.Fatal("expected error when destination exists without overwrite")
	}
}

func TestInherit_DestinationExists_Overwrite(t *testing.T) {
	st := makeStore(t)
	_ = st.Add("base", map[string]string{"A": "1"})
	_ = st.Add("child", map[string]string{"X": "9"})

	opts := inherit.DefaultOptions()
	opts.Overwrite = true
	if err := inherit.Profile(st, "base", "child", nil, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vars, _ := st.Get("child")
	if _, ok := vars["X"]; ok {
		t.Error("expected old key X to be gone after overwrite")
	}
}

func TestInherit_BaseNotFound_Strict(t *testing.T) {
	st := makeStore(t)
	err := inherit.Profile(st, "missing", "child", nil, inherit.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing base in strict mode")
	}
}

func TestInherit_BaseNotFound_NonStrict(t *testing.T) {
	st := makeStore(t)
	opts := inherit.DefaultOptions()
	opts.Strict = false
	if err := inherit.Profile(st, "missing", "child", map[string]string{"K": "v"}, opts); err != nil {
		t.Fatalf("unexpected error in non-strict mode: %v", err)
	}
	vars, _ := st.Get("child")
	if vars["K"] != "v" {
		t.Errorf("expected K=v from overrides, got %q", vars["K"])
	}
}

func TestInherit_EmptyDstName(t *testing.T) {
	st := makeStore(t)
	_ = st.Add("base", map[string]string{"A": "1"})
	err := inherit.Profile(st, "base", "", nil, inherit.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty destination name")
	}
}
