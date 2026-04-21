// Package prune removes variables from a profile that match given criteria,
// such as keys matching a pattern or values that are empty/blank.
package prune

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envoy-cli/internal/profile"
)

// DefaultOptions returns an Options with safe defaults.
func DefaultOptions() Options {
	return Options{
		DryRun:      false,
		EmptyValues: false,
	}
}

// Options controls prune behaviour.
type Options struct {
	// KeyPattern removes keys matching this regex (optional).
	KeyPattern string
	// EmptyValues removes keys whose value is empty or whitespace-only.
	EmptyValues bool
	// DryRun reports what would be removed without saving.
	DryRun bool
}

// Result holds the outcome of a prune operation.
type Result struct {
	Removed []string
	Kept    int
}

// Profile prunes variables from the named profile in the store and returns
// a Result describing what was (or would be) removed.
func Profile(store *profile.Store, name string, opts Options) (Result, error) {
	if name == "" {
		return Result{}, fmt.Errorf("profile name must not be empty")
	}

	vars, err := store.Get(name)
	if err != nil {
		return Result{}, fmt.Errorf("profile %q not found: %w", name, err)
	}

	var re *regexp.Regexp
	if opts.KeyPattern != "" {
		re, err = regexp.Compile(opts.KeyPattern)
		if err != nil {
			return Result{}, fmt.Errorf("invalid key pattern: %w", err)
		}
	}

	kept := make(map[string]string, len(vars))
	var removed []string

	for k, v := range vars {
		switch {
		case re != nil && re.MatchString(k):
			removed = append(removed, k)
		case opts.EmptyValues && strings.TrimSpace(v) == "":
			removed = append(removed, k)
		default:
			kept[k] = v
		}
	}

	if !opts.DryRun && len(removed) > 0 {
		if err := store.Add(name, kept); err != nil {
			return Result{}, fmt.Errorf("saving pruned profile: %w", err)
		}
	}

	return Result{Removed: removed, Kept: len(kept)}, nil
}
