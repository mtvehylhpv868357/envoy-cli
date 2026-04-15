package truncate_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/user/envoy-cli/internal/profile"
	"github.com/user/envoy-cli/internal/truncate"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "truncate-test-*")
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

func TestTruncate_Basic(t *testing.T) {
	st := makeStore(t)
	_ = st.Add("dev", map[string]string{"APP_HOST": "localhost", "DB_PASS": "secret", "PORT": "8080"})

	res, err := truncate.Profile(st, "dev", []string{"APP_HOST", "PORT"}, truncate.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sort.Strings(res.Kept)
	if len(res.Kept) != 2 || res.Kept[0] != "APP_HOST" || res.Kept[1] != "PORT" {
		t.Errorf("Kept = %v; want [APP_HOST PORT]", res.Kept)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "DB_PASS" {
		t.Errorf("Removed = %v; want [DB_PASS]", res.Removed)
	}

	vars, _ := st.Get("dev")
	if _, ok := vars["DB_PASS"]; ok {
		t.Error("DB_PASS should have been removed from saved profile")
	}
}

func TestTruncate_DryRun_DoesNotSave(t *testing.T) {
	st := makeStore(t)
	_ = st.Add("staging", map[string]string{"A": "1", "B": "2", "C": "3"})

	opts := truncate.DefaultOptions()
	opts.DryRun = true
	res, err := truncate.Profile(st, "staging", []string{"A"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed in dry-run, got %d", len(res.Removed))
	}

	vars, _ := st.Get("staging")
	if len(vars) != 3 {
		t.Errorf("dry-run should not modify store; got %d keys", len(vars))
	}
}

func TestTruncate_CaseInsensitive(t *testing.T) {
	st := makeStore(t)
	_ = st.Add("prod", map[string]string{"APP_HOST": "prod.example.com", "DB_URL": "postgres://..."})

	res, err := truncate.Profile(st, "prod", []string{"app_host"}, truncate.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Kept) != 1 {
		t.Errorf("expected 1 kept (case-insensitive), got %d", len(res.Kept))
	}
}

func TestTruncate_EmptyName_Error(t *testing.T) {
	st := makeStore(t)
	_, err := truncate.Profile(st, "", []string{"KEY"}, truncate.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty profile name")
	}
}

func TestTruncate_EmptyKeepKeys_Error(t *testing.T) {
	st := makeStore(t)
	_ = st.Add("dev", map[string]string{"A": "1"})
	_, err := truncate.Profile(st, "dev", []string{}, truncate.DefaultOptions())
	if err == nil {
		t.Error("expected error for empty keepKeys")
	}
}

func TestTruncate_ProfileNotFound_Error(t *testing.T) {
	st := makeStore(t)
	_, err := truncate.Profile(st, "ghost", []string{"KEY"}, truncate.DefaultOptions())
	if err == nil {
		t.Error("expected error for missing profile")
	}
}
