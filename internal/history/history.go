// Package history tracks which profiles have been activated and when.
package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single profile activation event.
type Entry struct {
	Profile   string    `json:"profile"`
	Activated time.Time `json:"activated"`
}

// Store manages the activation history log.
type Store struct {
	path string
}

// NewStore creates a Store backed by the given file path.
func NewStore(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return &Store{path: path}, nil
}

// Record appends a new activation entry for the given profile.
func (s *Store) Record(profile string) error {
	entries, err := s.ReadAll()
	if err != nil {
		return err
	}
	entries = append(entries, Entry{Profile: profile, Activated: time.Now().UTC()})
	return s.write(entries)
}

// ReadAll returns all recorded entries in chronological order.
func (s *Store) ReadAll() ([]Entry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// Last returns the most recently recorded entry, or nil if history is empty.
func (s *Store) Last() (*Entry, error) {
	entries, err := s.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	e := entries[len(entries)-1]
	return &e, nil
}

// Clear removes all history entries.
func (s *Store) Clear() error {
	return s.write([]Entry{})
}

func (s *Store) write(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
