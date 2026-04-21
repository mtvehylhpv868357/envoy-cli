// Package protect provides read-only locking for environment profiles,
// preventing accidental modification or deletion of critical profiles.
package protect

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const storeFile = "protected.json"

// Store holds the set of protected profile names.
type Store struct {
	path    string
	profile map[string]bool
}

// NewStore loads or initialises a protect store at dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Store{path: filepath.Join(dir, storeFile), profile: make(map[string]bool)}
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.profile); err != nil {
		return nil, err
	}
	return s, nil
}

// Protect marks a profile as protected.
func (s *Store) Protect(name string) error {
	if name == "" {
		return errors.New("protect: profile name must not be empty")
	}
	s.profile[name] = true
	return s.save()
}

// Unprotect removes the protection from a profile.
func (s *Store) Unprotect(name string) error {
	if name == "" {
		return errors.New("protect: profile name must not be empty")
	}
	delete(s.profile, name)
	return s.save()
}

// IsProtected reports whether the named profile is protected.
func (s *Store) IsProtected(name string) bool {
	return s.profile[name]
}

// List returns all currently protected profile names.
func (s *Store) List() []string {
	out := make([]string, 0, len(s.profile))
	for k := range s.profile {
		out = append(out, k)
	}
	return out
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.profile, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
