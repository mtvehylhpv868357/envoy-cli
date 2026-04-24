// Package envcount provides utilities for counting and summarising
// environment variable entries across one or more profiles.
package envcount

import (
	"fmt"
	"sort"
)

// Result holds the count statistics for a single profile.
type Result struct {
	Profile string
	Total   int
	Empty   int
	NonEmpty int
}

// String returns a human-readable summary line.
func (r Result) String() string {
	return fmt.Sprintf("%s: total=%d non-empty=%d empty=%d",
		r.Profile, r.Total, r.NonEmpty, r.Empty)
}

// Count returns a Result for the given profile name and its variable map.
func Count(profile string, vars map[string]string) Result {
	r := Result{Profile: profile, Total: len(vars)}
	for _, v := range vars {
		if v == "" {
			r.Empty++
		} else {
			r.NonEmpty++
		}
	}
	return r
}

// CountMany returns Results for each profile, sorted by profile name.
func CountMany(profiles map[string]map[string]string) []Result {
	names := make([]string, 0, len(profiles))
	for name := range profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	results := make([]Result, 0, len(names))
	for _, name := range names {
		results = append(results, Count(name, profiles[name]))
	}
	return results
}

// Total sums the variable counts across all provided Results.
func Total(results []Result) int {
	sum := 0
	for _, r := range results {
		sum += r.Total
	}
	return sum
}
