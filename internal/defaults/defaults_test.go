package defaults_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-cli/internal/defaults"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "defaults-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := filepath.Join(tempDir(t), "sub", "nested")
	_, err := defaults.NewStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestSet_And_Get(t *testing.T) {
	s, _ := defaults.NewStore(tempDir(t))
	if err := s.Set("APP_ENV", "development"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	v, ok := s.Get("APP_ENV")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if v != "development" {
		t.Errorf("got %q, want %q", v, "development")
	}
}

func TestSet_EmptyKey_Errors(t *testing.T) {
	s, _ := defaults.NewStore(tempDir(t))
	if err := s.Set("", "value"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestGet_NotFound(t *testing.T) {
	s, _ := defaults.NewStore(tempDir(t))
	_, ok := s.Get("MISSING")
	if ok {
		t.Error("expected key to be absent")
	}
}

func TestDelete_Existing(t *testing.T) {
	s, _ := defaults.NewStore(tempDir(t))
	_ = s.Set("LOG_LEVEL", "info")
	if err := s.Delete("LOG_LEVEL"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, ok := s.Get("LOG_LEVEL")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestDelete_NotFound_Errors(t *testing.T) {
	s, _ := defaults.NewStore(tempDir(t))
	if err := s.Delete("GHOST"); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s, _ := defaults.NewStore(tempDir(t))
	_ = s.Set("A", "1")
	_ = s.Set("B", "2")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
	// Mutating returned map should not affect store
	all["C"] = "3"
	if _, ok := s.Get("C"); ok {
		t.Error("mutation of returned map affected store")
	}
}

func TestApply_FillsMissing(t *testing.T) {
	s, _ := defaults.NewStore(tempDir(t))
	_ = s.Set("DB_HOST", "localhost")
	_ = s.Set("DB_PORT", "5432")
	vars := map[string]string{"DB_HOST": "prod.db"}
	result := s.Apply(vars)
	if result["DB_HOST"] != "prod.db" {
		t.Errorf("existing key overwritten: got %q", result["DB_HOST"])
	}
	if result["DB_PORT"] != "5432" {
		t.Errorf("missing key not filled: got %q", result["DB_PORT"])
	}
}
