// Package scope provides functionality for scoping environment variable
// profiles to specific directories or project paths.
package scope

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Binding associates a directory path with a named profile.
type Binding struct {
	Directory string `json:"directory"`
	Profile   string `json:"profile"`
}

// Store persists directory-to-profile bindings.
type Store struct {
	path     string
	bindings []Binding
}

// NewStore loads (or initialises) a scope store at the given path.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.bindings); err != nil {
		return nil, err
	}
	return s, nil
}

// Set binds a directory to a profile, replacing any existing binding.
func (s *Store) Set(dir, profile string) error {
	dir = filepath.Clean(dir)
	for i, b := range s.bindings {
		if b.Directory == dir {
			s.bindings[i].Profile = profile
			return s.save()
		}
	}
	s.bindings = append(s.bindings, Binding{Directory: dir, Profile: profile})
	return s.save()
}

// Get returns the profile bound to dir, or an empty string if none.
func (s *Store) Get(dir string) string {
	dir = filepath.Clean(dir)
	for _, b := range s.bindings {
		if b.Directory == dir {
			return b.Profile
		}
	}
	return ""
}

// Resolve walks up from dir toward the filesystem root, returning the
// first profile binding found, or an empty string.
func (s *Store) Resolve(dir string) string {
	dir = filepath.Clean(dir)
	for {
		if p := s.Get(dir); p != "" {
			return p
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// Remove deletes the binding for dir.
func (s *Store) Remove(dir string) error {
	dir = filepath.Clean(dir)
	filtered := s.bindings[:0]
	for _, b := range s.bindings {
		if b.Directory != dir {
			filtered = append(filtered, b)
		}
	}
	s.bindings = filtered
	return s.save()
}

// List returns all current bindings.
func (s *Store) List() []Binding {
	out := make([]Binding, len(s.bindings))
	copy(out, s.bindings)
	return out
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.bindings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
