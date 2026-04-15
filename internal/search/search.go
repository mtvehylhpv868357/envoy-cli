// Package search provides functionality to search for environment variables
// across profiles by key name, value pattern, or both.
package search

import (
	"strings"

	"github.com/user/envoy-cli/internal/profile"
)

// Result holds a single match found during a search.
type Result struct {
	Profile string
	Key     string
	Value   string
}

// Options controls how the search is performed.
type Options struct {
	// KeyPattern filters by key (case-insensitive substring).
	KeyPattern string
	// ValuePattern filters by value (case-insensitive substring).
	ValuePattern string
	// ExactKey requires an exact key match (case-insensitive).
	ExactKey bool
}

// DefaultOptions returns sensible search defaults.
func DefaultOptions() Options {
	return Options{}
}

// Profiles searches all profiles in the given store and returns matching results.
func Profiles(store *profile.Store, opts Options) ([]Result, error) {
	names := store.List()
	var results []Result

	for _, name := range names {
		p, err := store.Get(name)
		if err != nil {
			return nil, err
		}
		for k, v := range p.Vars {
			if !matchKey(k, opts) {
				continue
			}
			if !matchValue(v, opts) {
				continue
			}
			results = append(results, Result{
				Profile: name,
				Key:     k,
				Value:   v,
			})
		}
	}
	return results, nil
}

func matchKey(key string, opts Options) bool {
	if opts.KeyPattern == "" {
		return true
	}
	if opts.ExactKey {
		return strings.EqualFold(key, opts.KeyPattern)
	}
	return strings.Contains(strings.ToLower(key), strings.ToLower(opts.KeyPattern))
}

func matchValue(value string, opts Options) bool {
	if opts.ValuePattern == "" {
		return true
	}
	return strings.Contains(strings.ToLower(value), strings.ToLower(opts.ValuePattern))
}
