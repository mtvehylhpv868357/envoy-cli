// Package watch provides file-watching functionality for environment profile files.
// It monitors .env files for changes and triggers reload callbacks.
package watch

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Event represents a file change event.
type Event struct {
	Path    string
	ModTime time.Time
}

// Handler is called when a watched file changes.
type Handler func(Event)

// Watcher monitors a set of files for modifications.
type Watcher struct {
	mu       sync.Mutex
	files    map[string]time.Time
	handler  Handler
	interval time.Duration
	stop     chan struct{}
}

// New creates a new Watcher with the given poll interval and change handler.
func New(interval time.Duration, handler Handler) *Watcher {
	return &Watcher{
		files:    make(map[string]time.Time),
		handler:  handler,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Add registers a file path to be watched.
func (w *Watcher) Add(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("watch: resolving path %q: %w", path, err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return fmt.Errorf("watch: stat %q: %w", abs, err)
	}
	w.mu.Lock()
	w.files[abs] = info.ModTime()
	w.mu.Unlock()
	return nil
}

// Remove stops watching the given path.
func (w *Watcher) Remove(path string) {
	abs, _ := filepath.Abs(path)
	w.mu.Lock()
	delete(w.files, abs)
	w.mu.Unlock()
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.poll()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop halts the background polling goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) poll() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for path, last := range w.files {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.ModTime().After(last) {
			w.files[path] = info.ModTime()
			w.handler(Event{Path: path, ModTime: info.ModTime()})
		}
	}
}
