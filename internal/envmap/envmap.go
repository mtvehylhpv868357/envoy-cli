// Package envmap provides utilities for converting, merging, and
// manipulating environment variable maps.
package envmap

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// DefaultOptions returns a default set of options for envmap operations.
func DefaultOptions() Options {
	return Options{
		Overwrite: true,
		UppercaseKeys: false,
	}
}

// Options controls the behaviour of envmap operations.
type Options struct {
	Overwrite     bool
	UppercaseKeys bool
}

// FromEnviron converts a slice of "KEY=VALUE" strings (e.g. os.Environ())
// into a map.
func FromEnviron(environ []string) map[string]string {
	out := make(map[string]string, len(environ))
	for _, e := range environ {
		k, v, _ := strings.Cut(e, "=")
		if k != "" {
			out[k] = v
		}
	}
	return out
}

// ToEnviron converts a map into a sorted slice of "KEY=VALUE" strings
// suitable for use as a process environment.
func ToEnviron(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]string, 0, len(m))
	for _, k := range keys {
		out = append(out, fmt.Sprintf("%s=%s", k, m[k]))
	}
	return out
}

// Merge combines src into dst according to opts. Returns the merged map
// (dst is not modified).
func Merge(dst, src map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(dst)+len(src))
	for k, v := range dst {
		out[k] = v
	}
	for k, v := range src {
		key := k
		if opts.UppercaseKeys {
			key = strings.ToUpper(k)
		}
		if _, exists := out[key]; !exists || opts.Overwrite {
			out[key] = v
		}
	}
	return out
}

// FromOS returns a snapshot of the current process environment as a map.
func FromOS() map[string]string {
	return FromEnviron(os.Environ())
}

// Keys returns a sorted list of keys from m.
func Keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
