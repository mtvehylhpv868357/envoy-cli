// Package clone provides functionality to duplicate environment profiles
// under a new name, optionally overriding specific variables.
package clone

import (
	"fmt"

	"github.com/user/envoy-cli/internal/profile"
)

// Options controls the behaviour of a clone operation.
type Options struct {
	// Overrides are key=value pairs merged into the cloned profile.
	Overrides map[string]string
	// SetActive marks the new profile as the active one after cloning.
	SetActive bool
}

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Overrides: map[string]string{},
		SetActive: false,
	}
}

// Profile duplicates srcName into dstName within store, applying any
// overrides from opts. It returns an error if srcName does not exist or
// dstName already exists.
func Profile(store *profile.Store, srcName, dstName string, opts Options) error {
	if srcName == dstName {
.Errorf("clone: source and destination names must differ")
	}

	src, err := store.Get(srcName)
	if err != nil {
		return fmt.Errorf("clone: source profile %q not found: %w", srcName, err)
	}

	if _, err := store.Get(dstName); err == nil {
		return fmt.Errorf("clone: destination profile %q already exists", dstName)
	}

	// Deep-copy variables.
	vars := make(map[string]string, len(src.Vars))
	for k, v := range src.Vars {
		vars[k] = v
	}

	// Apply overrides.
	for k, v := range opts.Overrides {
		vars[k] = v
	}

	if err := store.Add(dstName, vars); err != nil {
		return fmt.Errorf("clone: could not create profile %q: %w", dstName, err)
	}

	if opts.SetActive {
		if err := store.SetActive(dstName); err != nil {
			return fmt.Errorf("clone: could not set %q as active: %w", dstName, err)
		}
	}

	return nil
}
