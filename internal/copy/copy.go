// Package copy provides functionality to duplicate environment variable profiles.
package copy

import (
	"errors"
	"fmt"

	"github.com/yourusername/envoy-cli/internal/profile"
)

// Options configures the copy operation.
type Options struct {
	// Overwrite allows the destination profile to be replaced if it already exists.
	Overwrite bool
}

// DefaultOptions returns the default copy options.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
	}
}

// Profile copies a profile from src to dst within the given store.
// If dst already exists and opts.Overwrite is false, an error is returned.
func Profile(store *profile.Store, src, dst string, opts Options) error {
	if src == "" {
		return errors.New("source profile name must not be empty")
	}
	if dst == "" {
		return errors.New("destination profile name must not be empty")
	}
	if src == dst {
		return errors.New("source and destination profile names must differ")
	}

	srcProfile, err := store.Get(src)
	if err != nil {
		return fmt.Errorf("source profile %q not found: %w", src, err)
	}

	if !opts.Overwrite {
		if _, err := store.Get(dst); err == nil {
			return fmt.Errorf("destination profile %q already exists; use --overwrite to replace it", dst)
		}
	}

	// Deep-copy the vars map so mutations don't affect the source.
	copied := make(map[string]string, len(srcProfile.Vars))
	for k, v := range srcProfile.Vars {
		copied[k] = v
	}

	if err := store.Add(dst, copied); err != nil {
		return fmt.Errorf("failed to save destination profile %q: %w", dst, err)
	}

	return nil
}
