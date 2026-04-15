package search_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/profile"
	"github.com/user/envoy-cli/internal/search"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "search-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

func makeStore(t *testing.T) *profile.Store {
	t.Helper()
	store, err := profile.LoadStore(filepath.Join(tempDir(t), "profiles"))
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func TestSearch_KeyPattern(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{"DB_HOST": "localhost", "APP_PORT": "8080"})
	_ = store.Add("prod", map[string]string{"DB_HOST": "prod-db", "SECRET": "xyz"})

	results, err := search.Profiles(store, search.Options{KeyPattern: "DB"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_ValuePattern(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{"DB_HOST": "localhost", "REDIS_HOST": "localhost"})
	_ = store.Add("prod", map[string]string{"DB_HOST": "prod-db"})

	results, err := search.Profiles(store, search.Options{ValuePattern: "localhost"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_ExactKey(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{"DB_HOST": "localhost", "DB_HOST_REPLICA": "replica"})

	results, err := search.Profiles(store, search.Options{KeyPattern: "DB_HOST", ExactKey: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", results[0].Key)
	}
}

func TestSearch_NoOptions_ReturnsAll(t *testing.T) {
	store := makeStore(t)
	_ = store.Add("dev", map[string]string{"A": "1", "B": "2"})

	results, err := search.Profiles(store, search.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearch_EmptyStore(t *testing.T) {
	store := makeStore(t)
	results, err := search.Profiles(store, search.Options{KeyPattern: "anything"})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
