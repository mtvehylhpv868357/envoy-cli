// Package promote provides functionality to promote (copy + activate) a
// profile from one named environment to another (e.g. staging → production).
package promote

import (
	"errors"
	"fmt"

	"github.com/yourusername/envoy-cli/internal/profile"
)

// Options controls the behaviour of a promotion.
type Options struct {
	// StorePath is the directory that holds profile JSON files.
	StorePath string
	// Overwrite allows the destination profile to be replaced if it already exists.
	Overwrite bool
	// Activate sets the promoted profile as the active profile after promotion.
	Activate bool
}

// DefaultOptions returns a sensible set of defaults.
func DefaultOptions() Options {
	return Options{
		Overwrite: false,
		Activate:  false,
	}
}

// Result describes what happened during a promotion.
type Result struct {
	Source      string
	Destination string
	Activated   bool
}

// Profile promotes the vars from src to dst inside the given store.
// If opts.Overwrite is false and dst already exists, an error is returned.
// If opts.Activate is true the destination profile becomes the active profile.
func Profile(src, dst string, opts Options) (Result, error) {
	if src == "" {
		return Result{}, errors.New("promote: source name must not be empty")
	}
	if dst == "" {
		return Result{}, errors.New("promote: destination name must not be empty")
	}
	if src == dst {
		return Result{}, fmt.Errorf("promote: source and destination are the same: %q", src)
	}

	store, err := profile.LoadStore(opts.StorePath)
	if err != nil {
		return Result{}, fmt.Errorf("promote: load store: %w", err)
	}

	srcVars, err := store.Get(src)
	if err != nil {
		return Result{}, fmt.Errorf("promote: source profile %q not found: %w", src, err)
	}

	if !opts.Overwrite {
		if _, err := store.Get(dst); err == nil {
			return Result{}, fmt.Errorf("promote: destination profile %q already exists (use --overwrite to replace)", dst)
		}
	}

	if err := store.Add(dst, srcVars); err != nil {
		return Result{}, fmt.Errorf("promote: write destination profile: %w", err)
	}

	result := Result{Source: src, Destination: dst}

	if opts.Activate {
		if err := store.SetActive(dst); err != nil {
			return result, fmt.Errorf("promote: set active: %w", err)
		}
		result.Activated = true
	}

	return result, nil
}
