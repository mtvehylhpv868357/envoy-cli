package profile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultStorageDir = ".envoy"
const profilesFile = "profiles.json"

// Profile represents a named set of environment variables.
type Profile struct {
	Name string            `json:"name"`
	Vars map[string]string `json:"vars"`
}

// Store holds all profiles for a project.
type Store struct {
	Profiles map[string]*Profile `json:"profiles"`
	Active   string              `json:"active"`
	path     string
}

// LoadStore reads the profiles store from disk, creating it if absent.
func LoadStore(baseDir string) (*Store, error) {
	dir := filepath.Join(baseDir, defaultStorageDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	path := filepath.Join(dir, profilesFile)
	store := &Store{Profiles: make(map[string]*Profile), path: path}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, store); err != nil {
		return nil, err
	}
	store.path = path
	return store, nil
}

// Save persists the store to disk.
func (s *Store) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Add creates or replaces a profile.
func (s *Store) Add(p *Profile) {
	s.Profiles[p.Name] = p
}

// Get returns a profile by name.
func (s *Store) Get(name string) (*Profile, bool) {
	p, ok := s.Profiles[name]
	return p, ok
}

// Delete removes a profile by name.
func (s *Store) Delete(name string) bool {
	if _, ok := s.Profiles[name]; !ok {
		return false
	}
	delete(s.Profiles, name)
	if s.Active == name {
		s.Active = ""
	}
	return true
}

// SetActive marks a profile as the active one.
func (s *Store) SetActive(name string) error {
	if _, ok := s.Profiles[name]; !ok {
		return errors.New("profile not found: " + name)
	}
	s.Active = name
	return nil
}
