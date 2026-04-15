// Package trim provides utilities for removing specific keys or values
// from an environment variable profile based on configurable criteria.
package trim

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/envoy-cli/internal/profile"
)

// DefaultOptions returns a sensible default Options value.
func DefaultOptions() Options {
	return Options{
		DryRun:      false,
		EmptyValues: false,
	}
}

// Options configures the Trim operation.
type Options struct {
	// DryRun reports what would be removed without modifying the profile.
	DryRun bool
	// Keys is an explicit list of keys to remove.
	Keys []string
	// Prefix removes all keys that start with the given prefix.
	Prefix string
	// EmptyValues removes all keys whose value is an empty string.
	EmptyValues bool
}

// Result holds the outcome of a Trim operation.
type Result struct {
	Removed []string
}

// Profile loads the named profile from store, removes matching keys
// according to opts, and (unless DryRun) saves the result back.
func Profile(store *profile.Store, name string, opts Options) (Result, error) {
	if name == "" {
		return Result{}, fmt.Errorf("profile name must not be empty")
	}

	p, err := store.Get(name)
	if err != nil {
		return Result{}, fmt.Errorf("profile %q not found: %w", name, err)
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	var removed []string
	updated := make(map[string]string, len(p.Vars))

	for k, v := range p.Vars {
		switch {
		case shouldRemove(k, v, keySet, opts):
			removed = append(removed, k)
		default:
			updated[k] = v
		}
	}

	if !opts.DryRun {
		p.Vars = updated
		if err := store.Save(p); err != nil {
			return Result{}, fmt.Errorf("saving profile: %w", err)
		}
	}

	return Result{Removed: removed}, nil
}

func shouldRemove(key, value string, keySet map[string]struct{}, opts Options) bool {
	if _, ok := keySet[key]; ok {
		return true
	}
	if opts.Prefix != "" && strings.HasPrefix(key, opts.Prefix) {
		return true
	}
	if opts.EmptyValues && value == "" {
		return true
	}
	return false
}
