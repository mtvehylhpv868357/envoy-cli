// Package stats provides aggregated statistics over environment variable profiles.
package stats

import (
	"sort"
	"strings"
)

// Report holds computed statistics for one or more profiles.
type Report struct {
	ProfileCount  int
	TotalVars     int
	UniqueKeys    int
	SharedKeys    int
	EmptyValues   int
	SensitiveKeys int
	KeyFrequency  map[string]int // how many profiles each key appears in
}

// sensitivePatterns are substrings that indicate a sensitive key.
var sensitivePatterns = []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE"}

// Compute derives a Report from a map of profileName -> vars.
func Compute(profiles map[string]map[string]string) Report {
	freq := make(map[string]int)

	totalVars := 0
	emptyValues := 0
	sensitiveKeys := 0

	for _, vars := range profiles {
		for k, v := range vars {
			freq[k]++
			totalVars++
			if v == "" {
				emptyValues++
			}
		}
	}

	// count unique vs shared
	uniqueKeys := 0
	sharedKeys := 0
	for _, count := range freq {
		if count == 1 {
			uniqueKeys++
		} else {
			sharedKeys++
		}
	}

	// count sensitive keys (distinct)
	for k := range freq {
		if isSensitive(k) {
			sensitiveKeys++
		}
	}

	return Report{
		ProfileCount:  len(profiles),
		TotalVars:     totalVars,
		UniqueKeys:    uniqueKeys,
		SharedKeys:    sharedKeys,
		EmptyValues:   emptyValues,
		SensitiveKeys: sensitiveKeys,
		KeyFrequency:  freq,
	}
}

// TopKeys returns the n most-common keys across all profiles, sorted by frequency desc.
func TopKeys(r Report, n int) []string {
	type kv struct {
		key   string
		count int
	}
	pairs := make([]kv, 0, len(r.KeyFrequency))
	for k, c := range r.KeyFrequency {
		pairs = append(pairs, kv{k, c})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].key < pairs[j].key
	})
	result := make([]string, 0, n)
	for i, p := range pairs {
		if i >= n {
			break
		}
		result = append(result, p.key)
	}
	return result
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}
