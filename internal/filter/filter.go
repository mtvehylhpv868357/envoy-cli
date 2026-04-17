package filter

import (
	"regexp"
	"strings"
)

// Options controls how filtering is applied.
type Options struct {
	KeyPattern   string
	ValuePattern string
	Prefix       string
	Invert       bool
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{}
}

// Map filters a map of env vars based on the provided options.
// Returns a new map containing only the matching entries.
func Map(vars map[string]string, opts Options) (map[string]string, error) {
	var keyRe, valRe *regexp.Regexp
	var err error

	if opts.KeyPattern != "" {
		if keyRe, err = regexp.Compile(opts.KeyPattern); err != nil {
			return nil, err
		}
	}
	if opts.ValuePattern != "" {
		if valRe, err = regexp.Compile(opts.ValuePattern); err != nil {
			return nil, err
		}
	}

	result := make(map[string]string)
	for k, v := range vars {
		matched := matchEntry(k, v, opts.Prefix, keyRe, valRe)
		if opts.Invert {
			matched = !matched
		}
		if matched {
			result[k] = v
		}
	}
	return result, nil
}

func matchEntry(key, value, prefix string, keyRe, valRe *regexp.Regexp) bool {
	if prefix != "" && !strings.HasPrefix(key, prefix) {
		return false
	}
	if keyRe != nil && !keyRe.MatchString(key) {
		return false
	}
	if valRe != nil && !valRe.MatchString(value) {
		return false
	}
	return true
}
