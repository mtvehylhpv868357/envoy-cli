// Package defaults provides functionality for managing default values
// for environment variable profiles. It allows users to define fallback
// values that are applied when a key is missing from the active profile.
package defaults

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultsFile = "defaults.json"

// Store manages a set of default key-value pairs persisted to disk.
type Store struct {
	path string
	data map[string]string
}

// NewStore opens or creates a defaults store at the given directory.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Store{path: filepath.Join(dir, defaultsFile), data: make(map[string]string)}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}

func (s *Store) save() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}

// Set stores a default value for the given key.
func (s *Store) Set(key, value string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	s.data[key] = value
	return s.save()
}

// Get returns the default value for a key and whether it exists.
func (s *Store) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

// Delete removes a default entry by key.
func (s *Store) Delete(key string) error {
	if _, ok := s.data[key]; !ok {
		return errors.New("key not found: " + key)
	}
	delete(s.data, key)
	return s.save()
}

// All returns a copy of all default key-value pairs.
func (s *Store) All() map[string]string {
	copy := make(map[string]string, len(s.data))
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}

// Apply merges defaults into vars, only setting keys that are absent.
func (s *Store) Apply(vars map[string]string) map[string]string {
	result := make(map[string]string, len(vars))
	for k, v := range vars {
		result[k] = v
	}
	for k, v := range s.data {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}
	return result
}
