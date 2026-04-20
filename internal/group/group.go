// Package group provides functionality for organizing profiles into named groups.
package group

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
)

// Store manages profile group memberships.
type Store struct {
	path   string
	groups map[string][]string // group name -> profile names
}

// NewStore loads or creates a group store at the given path.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path, groups: make(map[string][]string)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.groups); err != nil {
		return nil, err
	}
	return s, nil
}

// Add adds a profile to a group, deduplicating entries.
func (s *Store) Add(group, profile string) error {
	if group == "" {
		return errors.New("group name must not be empty")
	}
	if profile == "" {
		return errors.New("profile name must not be empty")
	}
	for _, p := range s.groups[group] {
		if p == profile {
			return nil
		}
	}
	s.groups[group] = append(s.groups[group], profile)
	return s.save()
}

// Remove removes a profile from a group.
func (s *Store) Remove(group, profile string) error {
	members := s.groups[group]
	filtered := members[:0]
	for _, p := range members {
		if p != profile {
			filtered = append(filtered, p)
		}
	}
	if len(filtered) == 0 {
		delete(s.groups, group)
	} else {
		s.groups[group] = filtered
	}
	return s.save()
}

// Get returns all profiles in a group.
func (s *Store) Get(group string) []string {
	return s.groups[group]
}

// List returns all group names sorted alphabetically.
func (s *Store) List() []string {
	names := make([]string, 0, len(s.groups))
	for k := range s.groups {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// Delete removes an entire group.
func (s *Store) Delete(group string) error {
	delete(s.groups, group)
	return s.save()
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.groups, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
