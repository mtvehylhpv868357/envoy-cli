// Package envlock provides functionality for locking environment variable
// profiles to prevent accidental modifications.
package envlock

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single lock record for a profile.
type Entry struct {
	Profile   string    `json:"profile"`
	LockedAt  time.Time `json:"locked_at"`
	Reason    string    `json:"reason,omitempty"`
}

// Store manages lock state for profiles.
type Store struct {
	dir string
}

// NewStore creates a new Store backed by the given directory.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("envlock: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) lockPath(profile string) string {
	return filepath.Join(s.dir, profile+".lock.json")
}

// Lock marks a profile as locked with an optional reason.
func (s *Store) Lock(profile, reason string) error {
	if profile == "" {
		return errors.New("envlock: profile name must not be empty")
	}
	entry := Entry{
		Profile:  profile,
		LockedAt: time.Now().UTC(),
		Reason:   reason,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("envlock: marshal: %w", err)
	}
	return os.WriteFile(s.lockPath(profile), data, 0o644)
}

// Unlock removes the lock for a profile.
func (s *Store) Unlock(profile string) error {
	if profile == "" {
		return errors.New("envlock: profile name must not be empty")
	}
	err := os.Remove(s.lockPath(profile))
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("envlock: profile %q is not locked", profile)
	}
	return err
}

// IsLocked reports whether a profile is currently locked.
func (s *Store) IsLocked(profile string) bool {
	_, err := os.Stat(s.lockPath(profile))
	return err == nil
}

// Get returns the lock entry for a profile, or an error if not locked.
func (s *Store) Get(profile string) (*Entry, error) {
	data, err := os.ReadFile(s.lockPath(profile))
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("envlock: profile %q is not locked", profile)
	}
	if err != nil {
		return nil, fmt.Errorf("envlock: read: %w", err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("envlock: unmarshal: %w", err)
	}
	return &entry, nil
}

// List returns all currently locked profiles.
func (s *Store) List() ([]Entry, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("envlock: read dir: %w", err)
	}
	var locks []Entry
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.dir, e.Name()))
		if err != nil {
			continue
		}
		var entry Entry
		if err := json.Unmarshal(data, &entry); err == nil {
			locks = append(locks, entry)
		}
	}
	return locks, nil
}
