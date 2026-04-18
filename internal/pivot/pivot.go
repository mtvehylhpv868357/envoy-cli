// Package pivot provides functionality to transpose environment variable
// profiles into a key-centric view across multiple profiles.
package pivot

import "sort"

// Row represents a single key's values across multiple profiles.
type Row struct {
	Key     string
	Values  map[string]string // profile name -> value
	Missing []string          // profile names where key is absent
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{SortKeys: true}
}

// Options controls pivot behaviour.
type Options struct {
	SortKeys bool
}

// Profiles pivots a map of profileName->vars into a slice of Rows,
// one per unique key across all profiles.
func Profiles(profiles map[string]map[string]string, opts Options) []Row {
	keySet := map[string]struct{}{}
	for _, vars := range profiles {
		for k := range vars {
			keySet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	if opts.SortKeys {
		sort.Strings(keys)
	}

	profileNames := make([]string, 0, len(profiles))
	for name := range profiles {
		profileNames = append(profileNames, name)
	}
	sort.Strings(profileNames)

	rows := make([]Row, 0, len(keys))
	for _, key := range keys {
		row := Row{
			Key:    key,
			Values: map[string]string{},
		}
		for _, pname := range profileNames {
			if v, ok := profiles[pname][key]; ok {
				row.Values[pname] = v
			} else {
				row.Missing = append(row.Missing, pname)
			}
		}
		rows = append(rows, row)
	}
	return rows
}
