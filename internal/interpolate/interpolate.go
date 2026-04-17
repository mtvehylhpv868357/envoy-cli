// Package interpolate provides variable interpolation for env var values.
// It supports ${VAR}, $VAR, and ${VAR:-default} syntax.
package interpolate

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)(?::-(.*?))?\}|\$([A-Z_][A-Z0-9_]*)`)

// Options controls interpolation behaviour.
type Options struct {
	// Strict causes an error if a referenced variable is not found and has no default.
	Strict bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Strict: false}
}

// Map interpolates all values in the provided map using other entries in the
// same map, falling back to os-style lookup via the lookup func.
func Map(vars map[string]string, lookup func(string) (string, bool), opts Options) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		resolved, err := interpolate(v, vars, lookup, opts)
		if err != nil {
			return nil, fmt.Errorf("key %s: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

// Value interpolates a single string value.
func Value(s string, vars map[string]string, lookup func(string) (string, bool), opts Options) (string, error) {
	return interpolate(s, vars, lookup, opts)
}

func interpolate(s string, vars map[string]string, lookup func(string) (string, bool), opts Options) (string, error) {
	var lastErr error
	result := refPattern.ReplaceAllStringFunc(s, func(match string) string {
		subs := refPattern.FindStringSubmatch(match)
		name := subs[1]
		defaultVal := subs[2]
		if name == "" {
			name = subs[3]
		}
		if v, ok := vars[name]; ok {
			return v
		}
		if lookup != nil {
			if v, ok := lookup(name); ok {
				return v
			}
		}
		if defaultVal != "" || strings.Contains(match, ":-") {
			return defaultVal
		}
		if opts.Strict {
			lastErr = fmt.Errorf("unresolved variable: %s", name)
		}
		return match
	})
	return result, lastErr
}
