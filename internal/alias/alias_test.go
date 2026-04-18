package alias

import (
	"os"
	"path/filepath"
	"testing"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "alias-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestNewStore_Empty(t *testing.T) {
	path := filepath.Join(tempDir(t), "aliases.json")
	s, err := NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.List()) != 0 {
		t.Error("expected empty store")
	}
}

func TestSet_And_Get(t *testing.T) {
	path := filepath.Join(tempDir(t), "aliases.json")
	s, _ := NewStore(path)
	if err := s.Set("prod", "production"); err != nil {
		t.Fatal(err)
	}
	v, ok := s.Get("prod")
	if !ok || v != "production" {
		t.Errorf("expected production, got %q", v)
	}
}

func TestSet_EmptyAlias_Errors(t *testing.T) {
	path := filepath.Join(tempDir(t), "aliases.json")
	s, _ := NewStore(path)
	if err := s.Set("", "production"); err == nil {
		t.Error("expected error for empty alias")
	}
}

func TestRemove_Existing(t *testing.T) {
	path := filepath.Join(tempDir(t), "aliases.json")
	s, _ := NewStore(path)
	_ = s.Set("dev", "development")
	if err := s.Remove("dev"); err != nil {
		t.Fatal(err)
	}
	if _, ok := s.Get("dev"); ok {
		t.Error("expected alias to be removed")
	}
}

func TestRemove_NotFound_Errors(t *testing.T) {
	path := filepath.Join(tempDir(t), "aliases.json")
	s, _ := NewStore(path)
	if err := s.Remove("ghost"); err == nil {
		t.Error("expected error for missing alias")
	}
}

func TestPersistence_RoundTrip(t *testing.T) {
	path := filepath.Join(tempDir(t), "aliases.json")
	s1, _ := NewStore(path)
	_ = s1.Set("stg", "staging")
	s2, err := NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	v, ok := s2.Get("stg")
	if !ok || v != "staging" {
		t.Errorf("round-trip failed, got %q", v)
	}
}
