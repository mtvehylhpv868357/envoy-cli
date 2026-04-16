// Package normalize provides utilities for normalizing environment variable maps.
package normalize

import (
	"strings"
)

// DefaultOptions returns a sensible default Options.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:  true,
		TrimSpace:      true,
		ReplaceHyphens: true,
	}
}

// Options controls normalization behaviour.
type Options struct {
	UppercaseKeys  bool
	TrimSpace      bool
	ReplaceHyphens bool // replace hyphens in keys with underscores
}

// Map applies normalization rules to a copy of the provided env map and
// returns the result together with a list of keys that were changed.
func Map(env map[string]string, opts Options) (map[string]string, []string) {
	out := make(map[string]string, len(env))
	var changed []string

	for k, v := range env {
		origKey := k
		origVal := v

		if opts.TrimSpace {
			k = strings.TrimSpace(k)
			v = strings.TrimSpace(v)
		}
		if opts.ReplaceHyphens {
			k = strings.ReplaceAll(k, "-", "_")
		}
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}

		if k != origKey || v != origVal {
			changed = append(changed, origKey)
		}
		out[k] = v
	}

	return out, changed
}

// Key normalizes a single key according to opts.
func Key(k string, opts Options) string {
	if opts.TrimSpace {
		k = strings.TrimSpace(k)
	}
	if opts.ReplaceHyphens {
		k = strings.ReplaceAll(k, "-", "_")
	}
	if opts.UppercaseKeys {
		k = strings.ToUpper(k)
	}
	return k
}
