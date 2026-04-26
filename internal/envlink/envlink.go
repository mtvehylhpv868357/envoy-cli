// Package envlink provides functionality for creating symbolic links between
// environment variable profiles, allowing one profile to reference another
// as its source of truth.
package envlink

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Link represents a symbolic reference from one profile to another.
type Link struct {
	Source  string    `json:"source"`
	Target  string    `json:"target"`
	Comment string    `json:"comment,omitempty"`
	Created time.Time `json:"created"`
}

// Store manages profile links persisted to disk.
type Store struct {
	dir string
}

// NewStore creates a new Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("envlink: create store dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) path(source string) string {
	return filepath.Join(s.dir, source+".link.json")
}

// Set creates or overwrites a link from source to target.
func (s *Store) Set(source, target, comment string) error {
	if source == "" {
		return errors.New("envlink: source name must not be empty")
	}
	if target == "" {
		return errors.New("envlink: target name must not be empty")
	}
	if source == target {
		return errors.New("envlink: source and target must differ")
	}
	link := Link{
		Source:  source,
		Target:  target,
		Comment: comment,
		Created: time.Now().UTC(),
	}
	data, err := json.MarshalIndent(link, "", "  ")
	if err != nil {
		return fmt.Errorf("envlink: marshal link: %w", err)
	}
	if err := os.WriteFile(s.path(source), data, 0o644); err != nil {
		return fmt.Errorf("envlink: write link: %w", err)
	}
	return nil
}

// Get retrieves the link for the given source profile.
func (s *Store) Get(source string) (Link, error) {
	data, err := os.ReadFile(s.path(source))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Link{}, fmt.Errorf("envlink: no link for %q", source)
		}
		return Link{}, fmt.Errorf("envlink: read link: %w", err)
	}
	var link Link
	if err := json.Unmarshal(data, &link); err != nil {
		return Link{}, fmt.Errorf("envlink: parse link: %w", err)
	}
	return link, nil
}

// Remove deletes the link for the given source profile.
func (s *Store) Remove(source string) error {
	err := os.Remove(s.path(source))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("envlink: remove link: %w", err)
	}
	return nil
}

// List returns all links stored in the store.
func (s *Store) List() ([]Link, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("envlink: list links: %w", err)
	}
	var links []Link
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.dir, e.Name()))
		if err != nil {
			continue
		}
		var link Link
		if err := json.Unmarshal(data, &link); err != nil {
			continue
		}
		links = append(links, link)
	}
	return links, nil
}

// IsLinked reports whether the given source profile has a link defined.
func (s *Store) IsLinked(source string) bool {
	_, err := os.Stat(s.path(source))
	return err == nil
}
