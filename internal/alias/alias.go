package alias

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Store manages short aliases that map to profile names.
type Store struct {
	path    string
	aliases map[string]string
}

// NewStore loads or creates an alias store at the given path.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path, aliases: make(map[string]string)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.aliases); err != nil {
		return nil, err
	}
	return s, nil
}

// Set creates or updates an alias pointing to a profile name.
func (s *Store) Set(alias, profile string) error {
	if alias == "" || profile == "" {
		return errors.New("alias and profile must not be empty")
	}
	s.aliases[alias] = profile
	return s.save()
}

// Get resolves an alias to a profile name.
func (s *Store) Get(alias string) (string, bool) {
	v, ok := s.aliases[alias]
	return v, ok
}

// Remove deletes an alias.
func (s *Store) Remove(alias string) error {
	if _, ok := s.aliases[alias]; !ok {
		return errors.New("alias not found: " + alias)
	}
	delete(s.aliases, alias)
	return s.save()
}

// List returns all aliases as a map copy.
func (s *Store) List() map[string]string {
	out := make(map[string]string, len(s.aliases))
	for k, v := range s.aliases {
		out[k] = v
	}
	return out
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.aliases, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
