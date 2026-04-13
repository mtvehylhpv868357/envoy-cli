package history_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/history"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "history-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNewStore(t *testing.T) {
	dir := tempDir(t)
	_, err := history.NewStore(filepath.Join(dir, "history.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
}

func TestReadAll_Empty(t *testing.T) {
	dir := tempDir(t)
	s, _ := history.NewStore(filepath.Join(dir, "history.json"))
	entries, err := s.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecord_AndReadAll(t *testing.T) {
	dir := tempDir(t)
	s, _ := history.NewStore(filepath.Join(dir, "history.json"))

	for _, p := range []string{"dev", "staging", "prod"} {
		if err := s.Record(p); err != nil {
			t.Fatalf("Record(%q): %v", p, err)
		}
	}

	entries, err := s.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[2].Profile != "prod" {
		t.Errorf("expected last profile to be 'prod', got %q", entries[2].Profile)
	}
}

func TestLast_NoEntries(t *testing.T) {
	dir := tempDir(t)
	s, _ := history.NewStore(filepath.Join(dir, "history.json"))
	e, err := s.Last()
	if err != nil {
		t.Fatalf("Last: %v", err)
	}
	if e != nil {
		t.Errorf("expected nil, got %+v", e)
	}
}

func TestLast_ReturnsLatest(t *testing.T) {
	dir := tempDir(t)
	s, _ := history.NewStore(filepath.Join(dir, "history.json"))
	_ = s.Record("alpha")
	_ = s.Record("beta")
	e, err := s.Last()
	if err != nil {
		t.Fatalf("Last: %v", err)
	}
	if e == nil || e.Profile != "beta" {
		t.Errorf("expected 'beta', got %v", e)
	}
}

func TestClear(t *testing.T) {
	dir := tempDir(t)
	s, _ := history.NewStore(filepath.Join(dir, "history.json"))
	_ = s.Record("dev")
	if err := s.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	entries, _ := s.ReadAll()
	if len(entries) != 0 {
		t.Errorf("expected empty after Clear, got %d entries", len(entries))
	}
}
