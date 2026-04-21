// Package envdiff provides utilities for comparing environment variable maps
// against the current OS environment, reporting missing, extra, and changed keys.
package envdiff

import (
	"os"
	"sort"
)

// Status represents the relationship of a key between two environments.
type Status string

const (
	StatusMissing  Status = "missing"  // present in profile, absent in OS env
	StatusExtra    Status = "extra"    // present in OS env, absent in profile
	StatusChanged  Status = "changed"  // present in both but values differ
	StatusMatching Status = "matching" // present in both with equal values
)

// Entry describes a single key's diff result.
type Entry struct {
	Key        string
	Status     Status
	ProfileVal string
	EnvVal     string
}

// Result holds the full diff output.
type Result struct {
	Entries []Entry
}

// Missing returns entries with StatusMissing.
func (r Result) Missing() []Entry { return r.filter(StatusMissing) }

// Extra returns entries with StatusExtra.
func (r Result) Extra() []Entry { return r.filter(StatusExtra) }

// Changed returns entries with StatusChanged.
func (r Result) Changed() []Entry { return r.filter(StatusChanged) }

func (r Result) filter(s Status) []Entry {
	out := []Entry{}
	for _, e := range r.Entries {
		if e.Status == s {
			out = append(out, e)
		}
	}
	return out
}

// CompareWithOS diffs profileVars against the live OS environment.
// Only keys present in profileVars are considered; extra OS keys are ignored
// unless includeExtra is true.
func CompareWithOS(profileVars map[string]string, includeExtra bool) Result {
	osEnv := osenviron()
	return compare(profileVars, osEnv, includeExtra)
}

// Compare diffs two explicit maps. If includeExtra is true, keys only in
// envB are included as StatusExtra entries.
func Compare(profileVars, envB map[string]string, includeExtra bool) Result {
	return compare(profileVars, envB, includeExtra)
}

func compare(a, b map[string]string, includeExtra bool) Result {
	seen := map[string]bool{}
	var entries []Entry

	for k, av := range a {
		seen[k] = true
		bv, ok := b[k]
		switch {
		case !ok:
			entries = append(entries, Entry{Key: k, Status: StatusMissing, ProfileVal: av})
		case av != bv:
			entries = append(entries, Entry{Key: k, Status: StatusChanged, ProfileVal: av, EnvVal: bv})
		default:
			entries = append(entries, Entry{Key: k, Status: StatusMatching, ProfileVal: av, EnvVal: bv})
		}
	}

	if includeExtra {
		for k, bv := range b {
			if !seen[k] {
				entries = append(entries, Entry{Key: k, Status: StatusExtra, EnvVal: bv})
			}
		}
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
	return Result{Entries: entries}
}

func osenviron() map[string]string {
	m := make(map[string]string)
	for _, e := range os.Environ() {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				m[e[:i]] = e[i+1:]
				break
			}
		}
	}
	return m
}
