// Package rollback provides functionality to revert a profile to a previous state
// using entries recorded in the history store.
package rollback

import (
	"errors"
	"fmt"
	"time"

	"github.com/yourusername/envoy-cli/internal/history"
	"github.com/yourusername/envoy-cli/internal/profile"
)

// ErrNoHistory is returned when there are no history entries for the given profile.
var ErrNoHistory = errors.New("no history entries found for profile")

// ErrIndexOutOfRange is returned when the requested index exceeds available entries.
var ErrIndexOutOfRange = errors.New("history index out of range")

// Result describes the outcome of a rollback operation.
type Result struct {
	Profile   string
	RolledTo  time.Time
	VarsCount int
}

// Options configures a rollback operation.
type Options struct {
	// Index is the 0-based position in the history list (0 = most recent).
	Index int
}

// DefaultOptions returns sensible defaults: roll back to the most recent entry.
func DefaultOptions() Options {
	return Options{Index: 0}
}

// Profile rolls back the named profile to a prior state recorded in history.
// It loads the history entry at opts.Index, restores the vars into the profile
// store, and returns a Result describing what was done.
func Profile(
	name string,
	store *profile.Store,
	hist *history.Store,
	opts Options,
) (Result, error) {
	if name == "" {
		return Result{}, errors.New("profile name must not be empty")
	}

	entries, err := hist.ReadAll(name)
	if err != nil {
		return Result{}, fmt.Errorf("reading history: %w", err)
	}
	if len(entries) == 0 {
		return Result{}, ErrNoHistory
	}
	if opts.Index < 0 || opts.Index >= len(entries) {
		return Result{}, fmt.Errorf("%w: requested %d, have %d", ErrIndexOutOfRange, opts.Index, len(entries))
	}

	// History is stored newest-first; entry[0] is the most recent snapshot.
	target := entries[opts.Index]

	if err := store.Set(name, target.Vars); err != nil {
		return Result{}, fmt.Errorf("restoring profile: %w", err)
	}

	return Result{
		Profile:   name,
		RolledTo:  target.Timestamp,
		VarsCount: len(target.Vars),
	}, nil
}
