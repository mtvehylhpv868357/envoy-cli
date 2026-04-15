// Package format provides utilities for rendering environment variable
// profiles as human-readable tables and structured output.
package format

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Row represents a single key-value pair with optional metadata.
type Row struct {
	Key    string
	Value  string
	Masked bool
}

// Options controls how a table is rendered.
type Options struct {
	// MaskSecrets replaces values of sensitive keys with asterisks.
	MaskSecrets bool
	// SecretPatterns is a list of substrings that mark a key as sensitive.
	SecretPatterns []string
	// NoColor disables ANSI colour codes.
	NoColor bool
}

// DefaultOptions returns sensible defaults for table rendering.
func DefaultOptions() Options {
	return Options{
		MaskSecrets:    true,
		SecretPatterns: []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS"},
	}
}

// Table writes env vars as an aligned two-column table to w.
func Table(w io.Writer, vars map[string]string, opts Options) {
	rows := buildRows(vars, opts)

	key // min width = len("KEY")
	for _, r := range rows {
		if len(r.Key) > keyWidth {
			keyWidth = len(r.Key)
		}
	}

	line := strings.Repeat("-", keyWidth+2+40)
	fmt.Fprintln(w, line)
	fmt.Fprintf(w, "%-*s  %s\n", keyWidth, "KEY", "VALUE")
	fmt.Fprintln(w, line)
	for _, r := range rows {
		v := r.Value
		if r.Masked {
			v = "********"
		}
		fmt.Fprintf(w, "%-*s  %s\n", keyWidth, r.Key, v)
	}
	fmt.Fprintln(w, line)
}

// buildRows converts a map into sorted, optionally masked rows.
func buildRows(vars map[string]string, opts Options) []Row {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rows := make([]Row, 0, len(keys))
	for _, k := range keys {
		masked := false
		if opts.MaskSecrets {
			for _, pat := range opts.SecretPatterns {
				if strings.Contains(strings.ToUpper(k), pat) {
					masked = true
					break
				}
			}
		}
		rows = append(rows, Row{Key: k, Value: vars[k], Masked: masked})
	}
	return rows
}
