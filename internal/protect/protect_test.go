package protect_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/yourorg/envoy-cli/internal/protect"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "protect-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNewStore_Empty(t *testing.T) {
	dir := tempDir(t)
	s, err := protect.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if got := s.List(); len(got) != 0 {
		t.Fatalf("expected empty list, got %v", got)
	}
}

func TestProtect_And_IsProtected(t *testing.T) {
	dir := tempDir(t)
	s, _ := protect.NewStore(dir)

	if err := s.Protect("production"); err != nil {
		t.Fatalf("Protect: %v", err)
	}
	if !s.IsProtected("production") {
		t.Fatal("expected production to be protected")
	}
	if s.IsProtected("staging") {
		t.Fatal("staging should not be protected")
	}
}

func TestUnprotect_RemovesEntry(t *testing.T) {
	dir := tempDir(t)
	s, _ := protect.NewStore(dir)
	_ = s.Protect("production")

	if err := s.Unprotect("production"); err != nil {
		t.Fatalf("Unprotect: %v", err)
	}
	if s.IsProtected("production") {
		t.Fatal("production should no longer be protected")
	}
}

func TestList_ReturnsAllProtected(t *testing.T) {
	dir := tempDir(t)
	s, _ := protect.NewStore(dir)
	_ = s.Protect("prod")
	_ = s.Protect("staging")

	list := s.List()
	slices.Sort(list)
	if len(list) != 2 || list[0] != "prod" || list[1] != "staging" {
		t.Fatalf("unexpected list: %v", list)
	}
}

func TestPersistence_AcrossStores(t *testing.T) {
	dir := tempDir(t)
	s1, _ := protect.NewStore(dir)
	_ = s1.Protect("prod")

	s2, err := protect.NewStore(dir)
	if err != nil {
		t.Fatalf("second NewStore: %v", err)
	}
	if !s2.IsProtected("prod") {
		t.Fatal("expected prod to be protected after reload")
	}
}

func TestProtect_EmptyName_Errors(t *testing.T) {
	dir := tempDir(t)
	s, _ := protect.NewStore(dir)
	if err := s.Protect(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestStoreFile_CreatedInDir(t *testing.T) {
	dir := tempDir(t)
	s, _ := protect.NewStore(dir)
	_ = s.Protect("prod")

	if _, err := os.Stat(filepath.Join(dir, "protected.json")); err != nil {
		t.Fatalf("expected store file to exist: %v", err)
	}
}
