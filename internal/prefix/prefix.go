// Package prefix provides utilities for adding or stripping key prefixes
// across environment variable profiles.
package prefix

import (
	"fmt"
	"strings"
)

// DefaultOptions returns sensible defaults for prefix operations.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
	}
}

// Options controls prefix behaviour.
type Options struct {
	// Overwrite allows existing keys to be replaced when a conflict arises.
	Overwrite bool
}

// Add returns a new map where every key is prefixed with p.
// If Overwrite is false and a destination key already exists the original
// (un-prefixed) entry is kept and an error is collected.
func Add(vars map[string]string, p string, opts Options) (map[string]string, []error) {
	out := make(map[string]string, len(vars))
	var errs []error

	for k, v := range vars {
		newKey := p + k
		if _, exists := out[newKey]; exists && !opts.Overwrite {
			errs = append(errs, fmt.Errorf("key conflict: %q already exists", newKey))
			continue
		}
		out[newKey] = v
	}
	return out, errs
}

// Strip returns a new map with the prefix p removed from every matching key.
// Keys that do not carry the prefix are passed through unchanged.
func Strip(vars map[string]string, p string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		stripped := strings.TrimPrefix(k, p)
		out[stripped] = v
	}
	return out
}

// FilterByPrefix returns only the entries whose keys start with p.
func FilterByPrefix(vars map[string]string, p string) map[string]string {
	out := make(map[string]string)
	for k, v := range vars {
		if strings.HasPrefix(k, p) {
			out[k] = v
		}
	}
	return out
}
