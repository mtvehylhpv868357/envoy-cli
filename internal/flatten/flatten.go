// Package flatten provides utilities for flattening nested environment
// variable maps using a delimiter-separated key convention.
package flatten

import (
	"fmt"
	"strings"
)

// DefaultOptions returns sensible defaults for flattening.
func DefaultOptions() Options {
	return Options{
		Delimiter: "_",
		Uppercase: true,
		Prefix:    "",
	}
}

// Options controls flattening behaviour.
type Options struct {
	Delimiter string
	Uppercase bool
	Prefix    string
}

// Map takes a nested map[string]any and returns a flat map[string]string.
// Nested keys are joined with the configured delimiter.
func Map(input map[string]any, opts Options) (map[string]string, error) {
	result := make(map[string]string)
	if err := flatten(input, opts.Prefix, opts, result); err != nil {
		return nil, err
	}
	return result, nil
}

func flatten(input map[string]any, prefix string, opts Options, out map[string]string) error {
	for k, v := range input {
		key := k
		if prefix != "" {
			key = prefix + opts.Delimiter + k
		}
		if opts.Uppercase {
			key = strings.ToUpper(key)
		}
		switch val := v.(type) {
		case map[string]any:
			if err := flatten(val, key, opts, out); err != nil {
				return err
			}
		case string:
			out[key] = val
		case nil:
			out[key] = ""
		default:
			out[key] = fmt.Sprintf("%v", val)
		}
	}
	return nil
}
