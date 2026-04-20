package group_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-cli/internal/group"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "group-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func TestNewStore_Empty(t *testing.T) {
	path := filepath.Join(tempDir(t), "groups.json")
	s, err := group.NewStore(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.List(); len(got) != 0 {
		t.Errorf("expected empty list, got %v", got)
	}
}

func TestAdd_And_Get(t *testing.T) {
	path := filepath.Join(tempDir(t), "groups.json")
	s, _ := group.NewStore(path)

	if err := s.Add("staging", "profile-a"); err != nil {
		t.Fatal(err)
	}
	if err := s.Add("staging", "profile-b"); err != nil {
		t.Fatal(err)
	}

	members := s.Get("staging")
	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

func TestAdd_Deduplicates(t *testing.T) {
	path := filepath.Join(tempDir(t), "groups.json")
	s, _ := group.NewStore(path)

	s.Add("g", "p1")
	s.Add("g", "p1")

	if got := len(s.Get("g")); got != 1 {
		t.Errorf("expected 1 member after dedup, got %d", got)
	}
}

func TestRemove_Profile(t *testing.T) {
	path := filepath.Join(tempDir(t), "groups.json")
	s, _ := group.NewStore(path)
	s.Add("g", "p1")
	s.Add("g", "p2")

	if err := s.Remove("g", "p1"); err != nil {
		t.Fatal(err)
	}
	members := s.Get("g")
	if len(members) != 1 || members[0] != "p2" {
		t.Errorf("unexpected members after remove: %v", members)
	}
}

func TestDelete_Group(t *testing.T) {
	path := filepath.Join(tempDir(t), "groups.json")
	s, _ := group.NewStore(path)
	s.Add("g", "p1")
	s.Delete("g")

	if got := s.List(); len(got) != 0 {
		t.Errorf("expected empty list after delete, got %v", got)
	}
}

func TestList_Sorted(t *testing.T) {
	path := filepath.Join(tempDir(t), "groups.json")
	s, _ := group.NewStore(path)
	s.Add("zebra", "p")
	s.Add("alpha", "p")
	s.Add("mango", "p")

	names := s.List()
	if names[0] != "alpha" || names[1] != "mango" || names[2] != "zebra" {
		t.Errorf("unexpected order: %v", names)
	}
}

func TestPersistence_RoundTrip(t *testing.T) {
	path := filepath.Join(tempDir(t), "groups.json")
	s1, _ := group.NewStore(path)
	s1.Add("prod", "profile-x")

	s2, err := group.NewStore(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if members := s2.Get("prod"); len(members) != 1 || members[0] != "profile-x" {
		t.Errorf("persistence failed: %v", members)
	}
}

func TestAdd_EmptyGroup_Errors(t *testing.T) {
	path := filepath.Join(tempDir(t), "groups.json")
	s, _ := group.NewStore(path)
	if err := s.Add("", "p"); err == nil {
		t.Error("expected error for empty group name")
	}
}
