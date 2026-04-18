// Package split provides functionality for splitting a profile into
// multiple profiles based on key prefixes or explicit key sets.
package split

import (
	"errors"
	"strings"

	"github.com/envoy-cli/internal/profile"
)

// DefaultOptions returns sensible defaults for Split.
func DefaultOptions() Options {
	return Options{
		StripPrefix: true,
		Overwrite:   false,
	}
}

// Options controls Split behaviour.
type Options struct {
	// StripPrefix removes the matched prefix from keys in destination profiles.
	StripPrefix bool
	// Overwrite allows overwriting existing destination profiles.
	Overwrite bool
}

// ByPrefix splits the source profile into one destination profile per prefix.
// Keys that match no prefix are collected into a profile named "default" when
// DefaultKey is non-empty, otherwise they are discarded.
func ByPrefix(store *profile.Store, source string, prefixes []string, opts Options) ([]string, error) {
	if source == "" {
		return nil, errors.New("split: source profile name must not be empty")
	}
	if len(prefixes) == 0 {
		return nil, errors.New("split: at least one prefix is required")
	}

	src, err := store.Get(source)
	if err != nil {
		return nil, err
	}

	buckets := make(map[string]map[string]string, len(prefixes))
	for _, p := range prefixes {
		buckets[p] = make(map[string]string)
	}

	for k, v := range src {
		for _, p := range prefixes {
			up := strings.ToUpper(p)
			if strings.HasPrefix(k, up) {
				destKey := k
				if opts.StripPrefix {
					destKey = strings.TrimPrefix(k, up)
				}
				if destKey != "" {
					buckets[p][destKey] = v
				}
				break
			}
		}
	}

	var created []string
	for _, p := range prefixes {
		vars := buckets[p]
		if len(vars) == 0 {
			continue
		}
		destName := strings.ToLower(strings.TrimSuffix(p, "_"))
		if !opts.Overwrite {
			if _, err := store.Get(destName); err == nil {
				return nil, errors.New("split: destination profile already exists: " + destName)
			}
		}
		if err := store.Add(destName, vars); err != nil {
			return nil, err
		}
		created = append(created, destName)
	}
	return created, nil
}
