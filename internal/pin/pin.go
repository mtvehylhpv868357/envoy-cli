// Package pin provides functionality for pinning environment variable profiles
// to specific directories, so that envoy-cli can auto-detect the active profile
// based on the current working directory.
package pin

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// ErrNotFound is returned when no pin exists for a given directory.
var ErrNotFound = errors.New("pin: no profile pinned to directory")

// Store persists directory→profile mappings.
type Store struct {
	path string
	pins map[string]string
}

// NewStore opens (or creates) a pin store at the given file path.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path, pins: make(map[string]string)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.pins); err != nil {
		return nil, err
	}
	return s, nil
}

// Set pins profileName to dir.
func (s *Store) Set(dir, profileName string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	s.pins[abs] = profileName
	return s.save()
}

// Get returns the profile pinned to dir, or ErrNotFound.
func (s *Store) Get(dir string) (string, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	name, ok := s.pins[abs]
	if !ok {
		return "", ErrNotFound
	}
	return name, nil
}

// Remove deletes the pin for dir.
func (s *Store) Remove(dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	if _, ok := s.pins[abs]; !ok {
		return ErrNotFound
	}
	delete(s.pins, abs)
	return s.save()
}

// List returns all pinned directory→profile pairs.
func (s *Store) List() map[string]string {
	out := make(map[string]string, len(s.pins))
	for k, v := range s.pins {
		out[k] = v
	}
	return out
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.pins, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
