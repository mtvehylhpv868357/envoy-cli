package expire_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/envoy-cli/envoy-cli/internal/expire"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "expire-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func newStore(t *testing.T) *expire.Store {
	t.Helper()
	path := filepath.Join(tempDir(t), "expire.json")
	s, err := expire.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestSet_And_Get(t *testing.T) {
	s := newStore(t)
	if err := s.Set("dev", time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}
	e, ok, err := s.Get("dev")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Profile != "dev" {
		t.Errorf("expected profile 'dev', got %q", e.Profile)
	}
	if time.Until(e.ExpiresAt) <= 0 {
		t.Error("expected expiry to be in the future")
	}
}

func TestGet_NotFound(t *testing.T) {
	s := newStore(t)
	_, ok, err := s.Get("missing")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if ok {
		t.Error("expected ok=false for missing profile")
	}
}

func TestIsExpired_Future(t *testing.T) {
	s := newStore(t)
	_ = s.Set("prod", time.Hour)
	expired, err := s.IsExpired("prod")
	if err != nil {
		t.Fatalf("IsExpired: %v", err)
	}
	if expired {
		t.Error("expected profile not to be expired")
	}
}

func TestIsExpired_Past(t *testing.T) {
	s := newStore(t)
	_ = s.Set("old", -time.Second)
	expired, err := s.IsExpired("old")
	if err != nil {
		t.Fatalf("IsExpired: %v", err)
	}
	if !expired {
		t.Error("expected profile to be expired")
	}
}

func TestRemove(t *testing.T) {
	s := newStore(t)
	_ = s.Set("staging", time.Hour)
	if err := s.Remove("staging"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	_, ok, _ := s.Get("staging")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestList(t *testing.T) {
	s := newStore(t)
	_ = s.Set("a", time.Hour)
	_ = s.Set("b", 2*time.Hour)
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestSet_EmptyName_Errors(t *testing.T) {
	s := newStore(t)
	if err := s.Set("", time.Hour); err == nil {
		t.Error("expected error for empty profile name")
	}
}
