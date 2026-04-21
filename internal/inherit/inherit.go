// Package inherit provides functionality for creating profiles that
// inherit variables from a base profile, with optional overrides.
package inherit

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/profile"
)

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
		Strict:    true,
	}
}

// Options controls the behaviour of Profile.
type Options struct {
	// Overwrite allows the destination profile to already exist.
	Overwrite bool
	// Strict causes an error when the base profile does not exist.
	Strict bool
}

// Profile creates a new profile named dst whose variables are the union of
// the base profile's variables and the provided overrides map.  Values in
// overrides take precedence over values inherited from base.
func Profile(st *profile.Store, base, dst string, overrides map[string]string, opts Options) error {
	if dst == "" {
		return fmt.Errorf("destination profile name must not be empty")
	}

	baseVars, err := st.Get(base)
	if err != nil {
		if opts.Strict {
			return fmt.Errorf("base profile %q not found: %w", base, err)
		}
		baseVars = map[string]string{}
	}

	if !opts.Overwrite {
		if _, err := st.Get(dst); err == nil {
			return fmt.Errorf("destination profile %q already exists; use overwrite option to replace", dst)
		}
	}

	merged := make(map[string]string, len(baseVars)+len(overrides))
	for k, v := range baseVars {
		merged[k] = v
	}
	for k, v := range overrides {
		merged[k] = v
	}

	return st.Add(dst, merged)
}
