package watch_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/envoy-cli/internal/watch"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "watch-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestAdd_NonExistentFile(t *testing.T) {
	w := watch.New(50*time.Millisecond, func(watch.Event) {})
	err := w.Add("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestAdd_ExistingFile(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("FOO=bar\n"), 0644); err != nil {
		t.Fatal(err)
	}
	w := watch.New(50*time.Millisecond, func(watch.Event) {})
	if err := w.Add(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("FOO=bar\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var mu sync.Mutex
	var got []watch.Event

	w := watch.New(20*time.Millisecond, func(e watch.Event) {
		mu.Lock()
		got = append(got, e)
		mu.Unlock()
	})
	if err := w.Add(path); err != nil {
		t.Fatal(err)
	}
	w.Start()
	defer w.Stop()

	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(path, []byte("FOO=changed\n"), 0644); err != nil {
		t.Fatal(err)
	}
	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	count := len(got)
	mu.Unlock()
	if count == 0 {
		t.Fatal("expected at least one change event, got none")
	}
	if got[0].Path != path {
		t.Errorf("event path = %q, want %q", got[0].Path, path)
	}
}

func TestWatcher_Remove(t *testing.T) {
	dir := tempDir(t)
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("A=1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var mu sync.Mutex
	var got []watch.Event

	w := watch.New(20*time.Millisecond, func(e watch.Event) {
		mu.Lock()
		got = append(got, e)
		mu.Unlock()
	})
	_ = w.Add(path)
	w.Remove(path)
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	_ = os.WriteFile(path, []byte("A=2\n"), 0644)
	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 0 {
		t.Errorf("expected no events after Remove, got %d", len(got))
	}
}
