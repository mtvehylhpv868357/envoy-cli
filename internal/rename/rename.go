// Package rename provides functionality for renaming environment profiles.
package rename

import (
	"errors"
	"fmt"

	"github.com/yourusername/envoy-cli/internal/profile"
)

// ErrSameName is returned when source and destination names are identical.
var ErrSameName = errors.New("source and destination names are the same")

// ErrNotFound is returned when the source profile does not exist.
var ErrNotFound = errors.New("profile not found")

// ErrAlreadyExists is returned when the destination profile already exists.
var ErrAlreadyExists = errors.New("destination profile already exists")

// Options configures the rename operation.
type Options struct {
	// Overwrite allows replacing an existing destination profile.
	Overwrite bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Overwrite: false}
}

// Profile renames a profile from oldName to newName within the given store.
func Profile(store *profile.Store, oldName, newName string, opts Options) error {
	if oldName == newName {
		return ErrSameName
	}

	vars, err := store.Get(oldName)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrNotFound, oldName)
	}

	if !opts.Overwrite {
		if _, err := store.Get(newName); err == nil {
			return fmt.Errorf("%w: %s", ErrAlreadyExists, newName)
		}
	}

	if err := store.Add(newName, vars); err != nil {
		return fmt.Errorf("creating destination profile: %w", err)
	}

	active, _ := store.ActiveName()
	if err := store.Delete(oldName); err != nil {
		return fmt.Errorf("removing old profile: %w", err)
	}

	// Preserve active profile pointer if it was the renamed one.
	if active == oldName {
		if err := store.SetActive(newName); err != nil {
			return fmt.Errorf("updating active profile: %w", err)
		}
	}

	return nil
}
