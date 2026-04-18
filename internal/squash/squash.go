// Package squash merges multiple profiles into a single flat map,
// applying them in order so later profiles override earlier ones.
package squash

import (
	"fmt"

	"github.com/envoy-cli/envoy-cli/internal/profile"
)

// DefaultOptions returns an Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Overwrite: true,
	}
}

// Options controls Squash behaviour.
type Options struct {
	// Overwrite allows later profiles to overwrite keys from earlier ones.
	Overwrite bool
}

// Profiles loads each named profile from store and merges their vars in order.
// If Overwrite is false, the first value for a key wins.
func Profiles(store *profile.Store, names []string, opts Options) (map[string]string, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("squash: at least one profile name required")
	}

	result := make(map[string]string)

	for _, name := range names {
		p, err := store.Get(name)
		if err != nil {
			return nil, fmt.Errorf("squash: profile %q not found: %w", name, err)
		}
		for k, v := range p.Vars {
			if _, exists := result[k]; exists && !opts.Overwrite {
				continue
			}
			result[k] = v
		}
	}

	return result, nil
}
