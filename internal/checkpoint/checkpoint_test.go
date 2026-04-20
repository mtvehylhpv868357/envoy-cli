package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/envoy-cli/internal/checkpoint"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "checkpoint-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := filepath.Join(tempDir(t), "nested", "checkpoints")
	_, err := checkpoint.NewStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestSave_And_Load(t *testing.T) {
	s, _ := checkpoint.NewStore(tempDir(t))
	e := checkpoint.Entry{
		Name:    "v1",
		Profile: "staging",
		Vars:    map[string]string{"FOO": "bar"},
		Note:    "before deploy",
	}
	if err := s.Save(e); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Load("v1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Profile != "staging" || got.Vars["FOO"] != "bar" || got.Note != "before deploy" {
		t.Errorf("unexpected entry: %+v", got)
	}
}

func TestSave_EmptyName_Errors(t *testing.T) {
	s, _ := checkpoint.NewStore(tempDir(t))
	if err := s.Save(checkpoint.Entry{}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestLoad_NotFound(t *testing.T) {
	s, _ := checkpoint.NewStore(tempDir(t))
	_, err := s.Load("missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDelete(t *testing.T) {
	s, _ := checkpoint.NewStore(tempDir(t))
	_ = s.Save(checkpoint.Entry{Name: "cp1", Profile: "dev", Vars: map[string]string{}})
	if err := s.Delete("cp1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Load("cp1")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestList_SortedByTime(t *testing.T) {
	s, _ := checkpoint.NewStore(tempDir(t))
	base := time.Now()
	_ = s.Save(checkpoint.Entry{Name: "b", Profile: "p", Vars: map[string]string{}, CreatedAt: base.Add(2 * time.Second)})
	_ = s.Save(checkpoint.Entry{Name: "a", Profile: "p", Vars: map[string]string{}, CreatedAt: base.Add(1 * time.Second)})
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Name != "a" || entries[1].Name != "b" {
		t.Errorf("wrong order: %v %v", entries[0].Name, entries[1].Name)
	}
}
