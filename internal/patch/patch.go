// Package patch applies partial updates to an existing profile.
package patch

import (
	"errors"
	"fmt"

	"github.com/envoy-cli/internal/profile"
)

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Delete: []string{},
		Upsert: map[string]string{},
	}
}

// Options controls patch behaviour.
type Options struct {
	Upsert map[string]string // keys to add or overwrite
	Delete []string          // keys to remove
}

// Profile applies the patch to the named profile in store.
func Profile(store *profile.Store, name string, opts Options) (map[string]string, error) {
	if name == "" {
		return nil, errors.New("patch: profile name must not be empty")
	}

	vars, err := store.Get(name)
	if err != nil {
		return nil, fmt.Errorf("patch: load profile %q: %w", name, err)
	}

	// copy so we don't mutate the original slice map
	result := make(map[string]string, len(vars))
	for k, v := range vars {
		result[k] = v
	}

	// apply deletions first
	for _, k := range opts.Delete {
		delete(result, k)
	}

	// apply upserts
	for k, v := range opts.Upsert {
		result[k] = v
	}

	if err := store.Add(name, result); err != nil {
		return nil, fmt.Errorf("patch: save profile %q: %w", name, err)
	}

	return result, nil
}
