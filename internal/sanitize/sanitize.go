// Package sanitize provides utilities for cleaning environment variable
// maps by normalizing keys, trimming whitespace, and removing invalid entries.
package sanitize

import (
	"strings"
	"unicode"
)

// Options controls sanitize behaviour.
type Options struct {
	// TrimSpace removes leading/trailing whitespace from values.
	TrimSpace bool
	// UppercaseKeys normalizes all keys to UPPER_CASE.
	UppercaseKeys bool
	// RemoveEmpty drops entries whose value is empty after trimming.
	RemoveEmpty bool
	// RemoveInvalidKeys drops entries whose key contains characters outside
	// [A-Za-z0-9_] or starts with a digit.
	RemoveInvalidKeys bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:         true,
		UppercaseKeys:     false,
		RemoveEmpty:       false,
		RemoveInvalidKeys: true,
	}
}

// Map applies the sanitize options to the provided env map and returns a new
// cleaned map. The original map is never mutated.
func Map(vars map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if opts.TrimSpace {
			v = strings.TrimSpace(v)
			k = strings.TrimSpace(k)
		}
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}
		if opts.RemoveInvalidKeys && !isValidKey(k) {
			continue
		}
		if opts.RemoveEmpty && v == "" {
			continue
		}
		out[k] = v
	}
	return out
}

// isValidKey returns true when k matches the POSIX env-var naming rules:
// starts with a letter or underscore, followed by letters, digits, or underscores.
func isValidKey(k string) bool {
	if k == "" {
		return false
	}
	for i, r := range k {
		switch {
		case i == 0 && (unicode.IsLetter(r) || r == '_'):
			// valid start
		case i > 0 && (unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'):
			// valid continuation
		default:
			return false
		}
	}
	return true
}
