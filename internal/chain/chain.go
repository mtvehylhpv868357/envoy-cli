// Package chain provides support for chaining multiple environment profiles
// together in a defined order, merging their variables with configurable
// precedence. Later profiles in the chain override earlier ones.
package chain

import (
	"errors"
	"fmt"

	"github.com/your-org/envoy-cli/internal/profile"
)

// DefaultOptions returns a new Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Overwrite: true,
		StopOnMissing: false,
	}
}

// Options controls how profiles are chained together.
type Options struct {
	// Overwrite determines whether later profiles override keys from earlier ones.
	// When false, the first value for a key wins.
	Overwrite bool

	// StopOnMissing causes Chain to return an error if any named profile does
	// not exist in the store, rather than silently skipping it.
	StopOnMissing bool
}

// Result holds the merged variables and metadata about the chain operation.
type Result struct {
	// Vars is the final merged map of environment variables.
	Vars map[string]string

	// Applied is the ordered list of profile names that were successfully merged.
	Applied []string

	// Skipped is the list of profile names that were missing and skipped
	// (only populated when StopOnMissing is false).
	Skipped []string
}

// Profiles merges the given profile names from the store in order.
// With default options, later profiles override earlier ones.
// Returns an error if StopOnMissing is true and a profile is not found.
func Profiles(store *profile.Store, names []string, opts Options) (*Result, error) {
	if len(names) == 0 {
		return nil, errors.New("chain: at least one profile name is required")
	}

	result := &Result{
		Vars:    make(map[string]string),
		Applied: make([]string, 0, len(names)),
		Skipped: []string{},
	}

	for _, name := range names {
		p, err := store.Get(name)
		if err != nil {
			if opts.StopOnMissing {
				return nil, fmt.Errorf("chain: profile %q not found", name)
			}
			result.Skipped = append(result.Skipped, name)
			continue
		}

		for k, v := range p.Vars {
			if _, exists := result.Vars[k]; exists && !opts.Overwrite {
				continue
			}
			result.Vars[k] = v
		}
		result.Applied = append(result.Applied, name)
	}

	if len(result.Applied) == 0 {
		return nil, errors.New("chain: no profiles could be applied")
	}

	return result, nil
}
