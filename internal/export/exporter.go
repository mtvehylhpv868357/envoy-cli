// Package export provides functionality for exporting environment profiles
// to various output formats such as .env files, shell scripts, and JSON.
package export

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Format represents the output format for exported variables.
type Format string

const (
	FormatDotEnv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Options configures the export behaviour.
type Options struct {
	Format  Format
	Sorted  bool
	Quote   bool
}

// DefaultOptions returns sensible export defaults.
func DefaultOptions() Options {
	return Options{
		Format: FormatDotEnv,
		Sorted: true,
		Quote:  true,
	}
}

// ToBytes serialises the provided env map to the requested format.
func ToBytes(vars map[string]string, opts Options) ([]byte, error) {
	switch opts.Format {
	case FormatDotEnv:
		return toDotEnv(vars, opts), nil
	case FormatExport:
		return toExport(vars, opts), nil
	case FormatJSON:
		return toJSON(vars)
	default:
		return nil, fmt.Errorf("export: unknown format %q", opts.Format)
	}
}

func keys(vars map[string]string, sorted bool) []string {
	ks := make([]string, 0, len(vars))
	for k := range vars {
		ks = append(ks, k)
	}
	if sorted {
		sort.Strings(ks)
	}
	return ks
}

func quoteValue(v string, quote bool) string {
	if !quote {
		return v
	}
	v = strings.ReplaceAll(v, `"`, `\"`)
	return `"` + v + `"`
}

func toDotEnv(vars map[string]string, opts Options) []byte {
	var sb strings.Builder
	for _, k := range keys(vars, opts.Sorted) {
		fmt.Fprintf(&sb, "%s=%s\n", k, quoteValue(vars[k], opts.Quote))
	}
	return []byte(sb.String())
}

func toExport(vars map[string]string, opts Options) []byte {
	var sb strings.Builder
	for _, k := range keys(vars, opts.Sorted) {
		fmt.Fprintf(&sb, "export %s=%s\n", k, quoteValue(vars[k], opts.Quote))
	}
	return []byte(sb.String())
}

func toJSON(vars map[string]string) ([]byte, error) {
	out, err := json.MarshalIndent(vars, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("export: json marshal: %w", err)
	}
	return append(out, '\n'), nil
}
