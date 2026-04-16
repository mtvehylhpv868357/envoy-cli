// Package inject provides utilities for injecting profile variables into
// a map or slice of strings (e.g. os.Environ style KEY=VALUE pairs).
package inject

import (
	"fmt"
	"strings"
)

// Options controls injection behaviour.
type Options struct {
	// Overwrite existing keys when true.
	Overwrite bool
	// Prefix is prepended to every injected key.
	Prefix string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Overwrite: true}
}

// IntoMap merges src into dst according to opts.
// dst is modified in place and also returned for convenience.
func IntoMap(dst, src map[string]string, opts Options) map[string]string {
	for k, v := range src {
		key := opts.Prefix + k
		if _, exists := dst[key]; exists && !opts.Overwrite {
			continue
		}
		dst[key] = v
	}
	return dst
}

// Intoenviron injects src into an os.Environ-style slice.
// Existing entries are updated or appended based on opts.
func IntoEnviron(environ []string, src map[string]string, opts Options) []string {
	// Build index of existing positions.
	index := make(map[string]int, len(environ))
	for i, entry := range environ {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			index[parts[0]] = i
		}
	}

	for k, v := range src {
		key := opts.Prefix + k
		pair := fmt.Sprintf("%s=%s", key, v)
		if i, exists := index[key]; exists {
			if opts.Overwrite {
				environ[i] = pair
			}
		} else {
			environ = append(environ, pair)
			index[key] = len(environ) - 1
		}
	}
	return environ
}
