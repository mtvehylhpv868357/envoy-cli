// Package diff provides utilities for comparing environment variable sets
// across profiles and snapshots.
package diff

import "sort"

// Change represents a single variable change between two env maps.
type Change struct {
	Key    string
	Old    string
	New    string
	Action string // "added", "removed", "modified"
}

// Compare returns the list of changes between two environment variable maps.
// `before` is the previous state, `after` is the new state.
func Compare(before, after map[string]string) []Change {
	var changes []Change

	// Check for removed or modified keys.
	for k, oldVal := range before {
		newVal, exists := after[k]
		if !exists {
			changes = append(changes, Change{Key: k, Old: oldVal, New: "", Action: "removed"})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Old: oldVal, New: newVal, Action: "modified"})
		}
	}

	// Check for added keys.
	for k, newVal := range after {
		if _, exists := before[k]; !exists {
			changes = append(changes, Change{Key: k, Old: "", New: newVal, Action: "added"})
		}
	}

	// Sort changes by key for deterministic output.
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return changes
}

// Summary returns a concise string summary of a slice of changes.
func Summary(changes []Change) map[string]int {
	summary := map[string]int{
		"added":    0,
		"removed":  0,
		"modified": 0,
	}
	for _, c := range changes {
		summary[c.Action]++
	}
	return summary
}
