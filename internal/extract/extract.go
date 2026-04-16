// Package extract provides functionality for extracting a subset of
// environment variables from a profile into a new profile.
package extract

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/envoy/internal/profile"
)

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
		Prefix:    "",
	}
}

// Options controls the behaviour of Profile.
type Options struct {
	// Overwrite allows the destination profile to be replaced if it already exists.
	Overwrite bool
	// Prefix filters keys to those that start with the given prefix (case-insensitive).
	// When empty all keys listed in Keys are used.
	Prefix string
	// Keys is an explicit list of variable names to extract.
	// Ignored when Prefix is non-empty.
	Keys []string
}

// Profile extracts a subset of variables from src into a new profile named dst.
// The store is used to load src and save the new dst profile.
func Profile(store *profile.Store, src, dst string, opts Options) (map[string]string, error) {
	if src == "" {
		return nil, fmt.Errorf("source profile name must not be empty")
	}
	if dst == "" {
		return nil, fmt.Errorf("destination profile name must not be empty")
	}
	if src == dst {
		return nil, fmt.Errorf("source and destination profile names must differ")
	}

	srcVars, err := store.Get(src)
	if err != nil {
		return nil, fmt.Errorf("load source profile %q: %w", src, err)
	}

	if !opts.Overwrite {
		if _, err := store.Get(dst); err == nil {
			return nil, fmt.Errorf("destination profile %q already exists; use Overwrite to replace it", dst)
		}
	}

	extracted := make(map[string]string)

	if opts.Prefix != "" {
		prefix := strings.ToUpper(opts.Prefix)
		for k, v := range srcVars {
			if strings.HasPrefix(strings.ToUpper(k), prefix) {
				extracted[k] = v
			}
		}
	} else {
		for _, k := range opts.Keys {
			v, ok := srcVars[k]
			if !ok {
				return nil, fmt.Errorf("key %q not found in source profile %q", k, src)
			}
			extracted[k] = v
		}
	}

	if len(extracted) == 0 {
		return nil, fmt.Errorf("no variables matched the extraction criteria")
	}

	if err := store.Add(dst, extracted); err != nil {
		return nil, fmt.Errorf("save destination profile %q: %w", dst, err)
	}

	return extracted, nil
}
