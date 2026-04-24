package envlock_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-cli/internal/envlock"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envlock-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := filepath.Join(tempDir(t), "locks")
	_, err := envlock.NewStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to exist: %v", err)
	}
}

func TestLock_And_IsLocked(t *testing.T) {
	s, _ := envlock.NewStore(tempDir(t))
	if err := s.Lock("production", "deploy freeze"); err != nil {
		t.Fatalf("Lock: %v", err)
	}
	if !s.IsLocked("production") {
		t.Error("expected production to be locked")
	}
	if s.IsLocked("staging") {
		t.Error("expected staging to be unlocked")
	}
}

func TestGet_ReturnsEntry(t *testing.T) {
	s, _ := envlock.NewStore(tempDir(t))
	_ = s.Lock("dev", "testing purposes")

	entry, err := s.Get("dev")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if entry.Profile != "dev" {
		t.Errorf("expected profile=dev, got %q", entry.Profile)
	}
	if entry.Reason != "testing purposes" {
		t.Errorf("expected reason='testing purposes', got %q", entry.Reason)
	}
	if entry.LockedAt.IsZero() {
		t.Error("expected non-zero LockedAt")
	}
}

func TestGet_NotFound(t *testing.T) {
	s, _ := envlock.NewStore(tempDir(t))
	_, err := s.Get("missing")
	if err == nil {
		t.Error("expected error for unlocked profile")
	}
}

func TestUnlock_RemovesLock(t *testing.T) {
	s, _ := envlock.NewStore(tempDir(t))
	_ = s.Lock("staging", "")
	if err := s.Unlock("staging"); err != nil {
		t.Fatalf("Unlock: %v", err)
	}
	if s.IsLocked("staging") {
		t.Error("expected staging to be unlocked after Unlock")
	}
}

func TestUnlock_NotLocked_Errors(t *testing.T) {
	s, _ := envlock.NewStore(tempDir(t))
	if err := s.Unlock("ghost"); err == nil {
		t.Error("expected error when unlocking non-locked profile")
	}
}

func TestList_ReturnsAllLocked(t *testing.T) {
	s, _ := envlock.NewStore(tempDir(t))
	_ = s.Lock("alpha", "")
	_ = s.Lock("beta", "reason")

	locks, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(locks) != 2 {
		t.Errorf("expected 2 locks, got %d", len(locks))
	}
}

func TestLock_EmptyName_Errors(t *testing.T) {
	s, _ := envlock.NewStore(tempDir(t))
	if err := s.Lock("", ""); err == nil {
		t.Error("expected error for empty profile name")
	}
}
