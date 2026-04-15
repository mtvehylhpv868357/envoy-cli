// Package merge provides utilities for merging environment variable profiles.
package merge

import (
	"fmt"

	"github.com/envoy-cli/internal/profile"
)

// Strategy defines how conflicting keys are resolved during a merge.
type Strategy string

const (
	// StrategyOurs keeps the value from the base profile on conflict.
	StrategyOurs Strategy = "ours"
	// StrategyTheirs keeps the value from the incoming profile on conflict.
	StrategyTheirs Strategy = "theirs"
)

// Options configures the merge behaviour.
type Options struct {
	Strategy  Strategy
	Overwrite bool // if true, save result back to base profile name
}

// DefaultOptions returns sensible merge defaults.
func DefaultOptions() Options {
	return Options{
		Strategy:  StrategyTheirs,
		Overwrite: false,
	}
}

// Result holds the outcome of a merge operation.
type Result struct {
	Merged   map[string]string
	Conflicts []string // keys that had conflicting values
}

// Profiles merges srcName into dstName using the provided store and options.
// The merged vars are returned in Result. If opts.Overwrite is true the
// destination profile is updated in the store.
func Profiles(store *profile.Store, dstName, srcName string, opts Options) (Result, error) {
	dst, err := store.Get(dstName)
	if err != nil {
		return Result{}, fmt.Errorf("merge: destination profile %q not found: %w", dstName, err)
	}

	src, err := store.Get(srcName)
	if err != nil {
		return Result{}, fmt.Errorf("merge: source profile %q not found: %w", srcName, err)
	}

	merged := make(map[string]string, len(dst.Vars))
	for k, v := range dst.Vars {
		merged[k] = v
	}

	var conflicts []string
	for k, v := range src.Vars {
		existing, exists := merged[k]
		if exists && existing != v {
			conflicts = append(conflicts, k)
			if opts.Strategy == StrategyTheirs {
				merged[k] = v
			}
		} else {
			merged[k] = v
		}
	}

	if opts.Overwrite {
		dst.Vars = merged
		if err := store.Add(dstName, dst.Vars); err != nil {
			return Result{}, fmt.Errorf("merge: failed to save profile: %w", err)
		}
	}

	return Result{Merged: merged, Conflicts: conflicts}, nil
}
