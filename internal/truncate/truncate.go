// Package truncate provides utilities for trimming environment variable
// profiles down to a subset of keys, optionally saving the result.
package truncate

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/profile"
)

// DefaultOptions returns a sensible Options value.
func DefaultOptions() Options {
	return Options{
		DryRun: false,
		CaseSensitive: false,
	}
}

// Options controls Truncate behaviour.
type Options struct {
	// DryRun prints what would be removed without modifying the profile.
	DryRun bool
	// CaseSensitive controls whether key matching is case-sensitive.
	CaseSensitive bool
}

// Result holds the outcome of a Truncate operation.
type Result struct {
	Kept    []string
	Removed []string
}

// Profile removes all keys from the named profile that are NOT present in
// keepKeys, then saves the modified profile back to the store.
func Profile(st *profile.Store, name string, keepKeys []string, opts Options) (Result, error) {
	if name == "" {
		return Result{}, fmt.Errorf("profile name must not be empty")
	}
	if len(keepKeys) == 0 {
		return Result{}, fmt.Errorf("keepKeys must contain at least one key")
	}

	vars, err := st.Get(name)
	if err != nil {
		return Result{}, fmt.Errorf("profile %q not found: %w", name, err)
	}

	keySet := make(map[string]struct{}, len(keepKeys))
	for _, k := range keepKeys {
		norm := k
		if !opts.CaseSensitive {
			norm = strings.ToUpper(k)
		}
		keySet[norm] = struct{}{}
	}

	result := Result{}
	newVars := make(map[string]string)

	for k, v := range vars {
		norm := k
		if !opts.CaseSensitive {
			norm = strings.ToUpper(k)
		}
		if _, ok := keySet[norm]; ok {
			newVars[k] = v
			result.Kept = append(result.Kept, k)
		} else {
			result.Removed = append(result.Removed, k)
		}
	}

	if !opts.DryRun {
		if err := st.Add(name, newVars); err != nil {
			return Result{}, fmt.Errorf("saving truncated profile: %w", err)
		}
	}

	return result, nil
}
