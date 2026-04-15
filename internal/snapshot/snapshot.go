package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot captures the state of an environment profile at a point in time.
type Snapshot struct {
	Name      string            `json:"name"`
	Profile   string            `json:"profile"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
}

// Store manages snapshots on disk.
type Store struct {
	dir string
}

// NewStore returns a Store rooted at dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save persists a snapshot to disk.
func (s *Store) Save(snap Snapshot) error {
	snap.CreatedAt = time.Now()
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := filepath.Join(s.dir, snap.Name+".json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write: %w", err)
	}
	return nil
}

// Load retrieves a snapshot by name.
func (s *Store) Load(name string) (*Snapshot, error) {
	path := filepath.Join(s.dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot %q not found", name)
		}
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}

// List returns the names of all stored snapshots.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("snapshot: list: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

// Delete removes a snapshot by name.
func (s *Store) Delete(name string) error {
	path := filepath.Join(s.dir, name+".json")
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("snapshot %q not found", name)
		}
		return fmt.Errorf("snapshot: delete: %w", err)
	}
	return nil
}

// Exists reports whether a snapshot with the given name exists in the store.
func (s *Store) Exists(name string) (bool, error) {
	path := filepath.Join(s.dir, name+".json")
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("snapshot: stat: %w", err)
}
