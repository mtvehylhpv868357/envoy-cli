// Package compare provides functionality to compare environment profiles
// side-by-side, showing differences between two named profiles.
package compare

import (
	"fmt"
	"sort"

	"github.com/user/envoy-cli/internal/profile"
)

// Result holds the comparison outcome between two profiles.
type Result struct {
	ProfileA string
	ProfileB string
	OnlyInA  map[string]string
	OnlyInB  map[string]string
	Differ   map[string][2]string // key -> [valueA, valueB]
	Same     map[string]string
}

// AllKeys returns a sorted list of all keys across both profiles.
func (r *Result) AllKeys() []string {
	seen := make(map[string]struct{})
	for k := range r.OnlyInA {
		seen[k] = struct{}{}
	}
	for k := range r.OnlyInB {
		seen[k] = struct{}{}
	}
	for k := range r.Differ {
		seen[k] = struct{}{}
	}
	for k := range r.Same {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Profiles compares two profiles by name using the provided store.
func Profiles(store *profile.Store, nameA, nameB string) (*Result, error) {
	pA, err := store.Get(nameA)
	if err != nil {
		return nil, fmt.Errorf("profile %q not found: %w", nameA, err)
	}
	pB, err := store.Get(nameB)
	if err != nil {
		return nil, fmt.Errorf("profile %q not found: %w", nameB, err)
	}

	result := &Result{
		ProfileA: nameA,
		ProfileB: nameB,
		OnlyInA:  make(map[string]string),
		OnlyInB:  make(map[string]string),
		Differ:   make(map[string][2]string),
		Same:     make(map[string]string),
	}

	for k, v := range pA.Vars {
		if vB, ok := pB.Vars[k]; !ok {
			result.OnlyInA[k] = v
		} else if v != vB {
			result.Differ[k] = [2]string{v, vB}
		} else {
			result.Same[k] = v
		}
	}
	for k, v := range pB.Vars {
		if _, ok := pA.Vars[k]; !ok {
			result.OnlyInB[k] = v
		}
	}
	return result, nil
}
