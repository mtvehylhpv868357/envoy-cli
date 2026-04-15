package tag_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/tag"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "tag-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestNewStore_CreatesEmpty(t *testing.T) {
	dir := tempDir(t)
	s, err := tag.NewStore(filepath.Join(dir, "tags.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get("dev"); len(got) != 0 {
		t.Errorf("expected empty tags, got %v", got)
	}
}

func TestAdd_And_Get(t *testing.T) {
	dir := tempDir(t)
	s, _ := tag.NewStore(filepath.Join(dir, "tags.json"))

	if err := s.Add("dev", "backend"); err != nil {
		t.Fatal(err)
	}
	if err := s.Add("dev", "local"); err != nil {
		t.Fatal(err)
	}

	tags := s.Get("dev")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "backend" || tags[1] != "local" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestAdd_Deduplicates(t *testing.T) {
	dir := tempDir(t)
	s, _ := tag.NewStore(filepath.Join(dir, "tags.json"))
	s.Add("dev", "backend")
	s.Add("dev", "backend")
	if len(s.Get("dev")) != 1 {
		t.Errorf("expected 1 tag after duplicate add")
	}
}

func TestRemove_Tag(t *testing.T) {
	dir := tempDir(t)
	s, _ := tag.NewStore(filepath.Join(dir, "tags.json"))
	s.Add("dev", "backend")
	s.Add("dev", "local")
	s.Remove("dev", "backend")
	tags := s.Get("dev")
	if len(tags) != 1 || tags[0] != "local" {
		t.Errorf("unexpected tags after remove: %v", tags)
	}
}

func TestFindByTag(t *testing.T) {
	dir := tempDir(t)
	s, _ := tag.NewStore(filepath.Join(dir, "tags.json"))
	s.Add("dev", "backend")
	s.Add("staging", "backend")
	s.Add("prod", "frontend")

	results := s.FindByTag("backend")
	if len(results) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(results))
	}
	if results[0] != "dev" || results[1] != "staging" {
		t.Errorf("unexpected profiles: %v", results)
	}
}

func TestPersistence(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, "tags.json")
	s1, _ := tag.NewStore(path)
	s1.Add("dev", "backend")

	s2, err := tag.NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if tags := s2.Get("dev"); len(tags) != 1 || tags[0] != "backend" {
		t.Errorf("tags not persisted: %v", tags)
	}
}
