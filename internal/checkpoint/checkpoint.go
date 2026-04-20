// Package checkpoint provides named save-points for environment profiles,
// allowing users to bookmark a profile state and restore it later.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a named checkpoint for a profile.
type Entry struct {
	Name      string            `json:"name"`
	Profile   string            `json:"profile"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
	Note      string            `json:"note,omitempty"`
}

// Store manages checkpoint persistence.
type Store struct {
	dir string
}

// NewStore returns a Store rooted at dir, creating it if necessary.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name+".json")
}

// Save writes a checkpoint entry to disk, overwriting any existing one with the same name.
func (s *Store) Save(e Entry) error {
	if e.Name == "" {
		return errors.New("checkpoint name must not be empty")
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(e.Name), data, 0o644)
}

// Load retrieves a checkpoint by name.
func (s *Store) Load(name string) (Entry, error) {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Entry{}, errors.New("checkpoint not found: " + name)
		}
		return Entry{}, err
	}
	var e Entry
	return e, json.Unmarshal(data, &e)
}

// Delete removes a checkpoint by name.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("checkpoint not found: " + name)
	}
	return err
}

// List returns all stored checkpoints sorted by creation time.
func (s *Store) List() ([]Entry, error) {
	glob, err := filepath.Glob(filepath.Join(s.dir, "*.json"))
	if err != nil {
		return nil, err
	}
	var entries []Entry
	for _, p := range glob {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err == nil {
			entries = append(entries, e)
		}
	}
	sortByTime(entries)
	return entries, nil
}

func sortByTime(entries []Entry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].CreatedAt.Before(entries[j-1].CreatedAt); j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
