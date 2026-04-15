// Package resolve provides variable interpolation for environment profiles.
// It expands references like ${VAR} or $VAR within values using a given map.
package resolve

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls resolution behaviour.
type Options struct {
	// FallbackToOS allows falling back to the process environment when a
	// variable is not found in the provided map.
	FallbackToOS bool
	// Strict causes Resolve to return an error if any reference is unresolved.
	Strict bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		FallbackToOS: true,
		Strict:       false,
	}
}

// Vars resolves all variable references within the values of the provided map
// and returns a new map with expanded values.
func Vars(vars map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		resolved, err := resolveValue(v, vars, opts)
		if err != nil {
			return nil, fmt.Errorf("resolve %s: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

// Value resolves variable references within a single string using the provided
// variable map and options.
func Value(s string, vars map[string]string, opts Options) (string, error) {
	return resolveValue(s, vars, opts)
}

func resolveValue(s string, vars map[string]string, opts Options) (string, error) {
	var firstErr error
	result := refPattern.ReplaceAllStringFunc(s, func(match string) string {
		name := extractName(match)
		if val, ok := vars[name]; ok {
			return val
		}
		if opts.FallbackToOS {
			if val, ok := os.LookupEnv(name); ok {
				return val
			}
		}
		if opts.Strict && firstErr == nil {
			firstErr = fmt.Errorf("unresolved variable: %s", name)
		}
		return match
	})
	if firstErr != nil {
		return "", firstErr
	}
	return result, nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "${") 
	match = strings.TrimSuffix(match, "}")
	match = strings.TrimPrefix(match, "$")
	return match
}
