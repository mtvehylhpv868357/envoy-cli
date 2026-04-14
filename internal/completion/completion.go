// Package completion provides shell auto-completion helpers for envoy-cli.
// It generates completion scripts for profile names, snapshot names, and
// other dynamic values used across commands.
package completion

import (
	"os"
	"path/filepath"
	"strings"
)

// ProfileNames returns a list of profile names found in the given store
// directory. Returns an empty slice if the directory does not exist.
func ProfileNames(storeDir string) []string {
	entries, err := os.ReadDir(storeDir)
	if err != nil {
		return []string{}
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".json") {
			names = append(names, strings.TrimSuffix(name, ".json"))
		}
	}
	return names
}

// SnapshotNames returns a list of snapshot names found in the given store
// directory. Each snapshot is stored as a <name>.json file.
func SnapshotNames(storeDir string) []string {
	return ProfileNames(storeDir) // same layout
}

// EnvFiles returns .env files found in dir (non-recursive).
func EnvFiles(dir string) []string {
	matches, err := filepath.Glob(filepath.Join(dir, ".env*"))
	if err != nil {
		return []string{}
	}
	return matches
}

// FilterPrefix returns only those items whose name starts with prefix.
func FilterPrefix(items []string, prefix string) []string {
	if prefix == "" {
		return items
	}
	var out []string
	for _, item := range items {
		if strings.HasPrefix(item, prefix) {
			out = append(out, item)
		}
	}
	return out
}
