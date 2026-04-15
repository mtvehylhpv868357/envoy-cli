package pin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-cli/internal/pin"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "pin-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestNewStore_Empty(t *testing.T) {
	dir := tempDir(t)
	s, err := pin.NewStore(filepath.Join(dir, "pins.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.List(); len(got) != 0 {
		t.Errorf("expected empty store, got %v", got)
	}
}

func TestSet_And_Get(t *testing.T) {
	dir := tempDir(t)
	s, _ := pin.NewStore(filepath.Join(dir, "pins.json"))

	projectDir := filepath.Join(dir, "myproject")
	if err := s.Set(projectDir, "production"); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, err := s.Get(projectDir)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "production" {
		t.Errorf("expected 'production', got %q", got)
	}
}

func TestGet_NotFound(t *testing.T) {
	dir := tempDir(t)
	s, _ := pin.NewStore(filepath.Join(dir, "pins.json"))

	_, err := s.Get(filepath.Join(dir, "nowhere"))
	if err != pin.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestRemove(t *testing.T) {
	dir := tempDir(t)
	s, _ := pin.NewStore(filepath.Join(dir, "pins.json"))
	projectDir := filepath.Join(dir, "proj")

	s.Set(projectDir, "staging")
	if err := s.Remove(projectDir); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, err := s.Get(projectDir)
	if err != pin.ErrNotFound {
		t.Errorf("expected ErrNotFound after remove, got %v", err)
	}
}

func TestRemove_NotFound(t *testing.T) {
	dir := tempDir(t)
	s, _ := pin.NewStore(filepath.Join(dir, "pins.json"))
	err := s.Remove(filepath.Join(dir, "ghost"))
	if err != pin.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestPersistence(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "pins.json")
	projectDir := filepath.Join(dir, "app")

	s1, _ := pin.NewStore(path)
	s1.Set(projectDir, "dev")

	s2, err := pin.NewStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	got, err := s2.Get(projectDir)
	if err != nil {
		t.Fatalf("Get after reload: %v", err)
	}
	if got != "dev" {
		t.Errorf("expected 'dev', got %q", got)
	}
}
