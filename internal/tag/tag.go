// Package tag provides functionality for tagging and filtering environment profiles.
package tag

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
)

// TagStore maps profile names to their associated tags.
type TagStore struct {
	path string
	data map[string][]string // profile -> tags
}

// NewStore opens or creates a tag store at the given path.
func NewStore(path string) (*TagStore, error) {
	s := &TagStore{path: path, data: make(map[string][]string)}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Add adds a tag to a profile (no-op if already present).
func (s *TagStore) Add(profile, tag string) error {
	for _, t := range s.data[profile] {
		if t == tag {
			return nil
		}
	}
	s.data[profile] = append(s.data[profile], tag)
	sort.Strings(s.data[profile])
	return s.save()
}

// Remove removes a tag from a profile.
func (s *TagStore) Remove(profile, tag string) error {
	tags := s.data[profile]
	updated := tags[:0]
	for _, t := range tags {
		if t != tag {
			updated = append(updated, t)
		}
	}
	s.data[profile] = updated
	return s.save()
}

// Get returns all tags for a profile.
func (s *TagStore) Get(profile string) []string {
	return s.data[profile]
}

// FindByTag returns all profiles that have the given tag.
func (s *TagStore) FindByTag(tag string) []string {
	var results []string
	for profile, tags := range s.data {
		for _, t := range tags {
			if t == tag {
				results = append(results, profile)
				break
			}
		}
	}
	sort.Strings(results)
	return results
}

func (s *TagStore) load() error {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}

func (s *TagStore) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}
