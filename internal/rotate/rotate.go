// Package rotate provides functionality to rotate (regenerate) environment variable values
// within a profile by applying a transformation function.
package rotate

import (
	"errors"
	"fmt"

	"github.com/user/envoy-cli/internal/profile"
)

// DefaultOptions returns a sensible default Options.
func DefaultOptions() Options {
	return Options{
		DryRun: false,
	}
}

// Options controls the behaviour of Profile.
type Options struct {
	// DryRun prevents any writes when true.
	DryRun bool
	// Keys limits rotation to the specified keys; all keys are rotated when empty.
	Keys []string
}

// TransformFunc receives the key and current value and returns the new value.
type TransformFunc func(key, value string) (string, error)

// Profile rotates values in the named profile using fn and saves the result.
// It returns a map of key -> new value for every key that was changed.
func Profile(store *profile.Store, name string, fn TransformFunc, opts Options) (map[string]string, error) {
	if name == "" {
		return nil, errors.New("rotate: profile name must not be empty")
	}
	if fn == nil {
		return nil, errors.New("rotate: transform function must not be nil")
	}

	p, err := store.Get(name)
	if err != nil {
		return nil, fmt.Errorf("rotate: %w", err)
	}

	target := buildTargetSet(opts.Keys)

	changed := make(map[string]string)
	for k, v := range p.Vars {
		if len(target) > 0 && !target[k] {
			continue
		}
		newVal, err := fn(k, v)
		if err != nil {
			return nil, fmt.Errorf("rotate: transform failed for key %q: %w", k, err)
		}
		if newVal != v {
			changed[k] = newVal
			p.Vars[k] = newVal
		}
	}

	if !opts.DryRun && len(changed) > 0 {
		if err := store.Add(name, p.Vars); err != nil {
			return nil, fmt.Errorf("rotate: save failed: %w", err)
		}
	}

	return changed, nil
}

func buildTargetSet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
