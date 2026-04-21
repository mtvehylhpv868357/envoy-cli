// Package expire provides TTL-based expiration tracking for environment profiles.
package expire

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Entry holds expiration metadata for a single profile.
type Entry struct {
	Profile   string    `json:"profile"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Store manages expiration entries persisted to disk.
type Store struct {
	path string
}

// NewStore returns a Store backed by the given file path.
func NewStore(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return &Store{path: path}, nil
}

func (s *Store) load() ([]Entry, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Entry
	return entries, json.Unmarshal(data, &entries)
}

func (s *Store) save(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Set records an expiration time for the given profile.
func (s *Store) Set(profile string, ttl time.Duration) error {
	if profile == "" {
		return errors.New("profile name must not be empty")
	}
	entries, err := s.load()
	if err != nil {
		return err
	}
	for i, e := range entries {
		if e.Profile == profile {
			entries[i].ExpiresAt = time.Now().Add(ttl)
			return s.save(entries)
		}
	}
	entries = append(entries, Entry{Profile: profile, ExpiresAt: time.Now().Add(ttl)})
	return s.save(entries)
}

// Get returns the expiration entry for a profile, or false if not found.
func (s *Store) Get(profile string) (Entry, bool, error) {
	entries, err := s.load()
	if err != nil {
		return Entry{}, false, err
	}
	for _, e := range entries {
		if e.Profile == profile {
			return e, true, nil
		}
	}
	return Entry{}, false, nil
}

// IsExpired reports whether the profile has a recorded expiration that has passed.
func (s *Store) IsExpired(profile string) (bool, error) {
	e, ok, err := s.Get(profile)
	if err != nil || !ok {
		return false, err
	}
	return time.Now().After(e.ExpiresAt), nil
}

// Remove deletes the expiration entry for a profile.
func (s *Store) Remove(profile string) error {
	entries, err := s.load()
	if err != nil {
		return err
	}
	filtered := entries[:0]
	for _, e := range entries {
		if e.Profile != profile {
			filtered = append(filtered, e)
		}
	}
	return s.save(filtered)
}

// List returns all expiration entries.
func (s *Store) List() ([]Entry, error) {
	return s.load()
}
