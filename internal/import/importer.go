// Package importer provides functionality for importing environment variables
// from various external formats into envoy-cli profiles.
package importer

import (
	"fmt"
	"os"
	"strings"
)

// Format represents a supported import format.
type Format string

const (
	FormatDotEnv  Format = "dotenv"
	FormatExport  Format = "export"
	FormatJSON    Format = "json"
)

// Result holds the outcome of an import operation.
type Result struct {
	Vars    map[string]string
	Skipped []string
	Format  Format
}

// FromFile reads a file and attempts to parse it using the given format.
// If format is empty, it is inferred from the file contents.
func FromFile(path string, format Format) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("import: cannot read file %q: %w", path, err)
	}
	return FromBytes(data, format)
}

// FromBytes parses raw bytes using the specified format.
func FromBytes(data []byte, format Format) (*Result, error) {
	if format == "" {
		format = detect(data)
	}
	switch format {
	case FormatDotEnv:
		return parseDotEnv(string(data))
	case FormatExport:
		return parseExport(string(data))
	default:
		return nil, fmt.Errorf("import: unsupported format %q", format)
	}
}

// detect guesses the format from file contents.
func detect(data []byte) Format {
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "export ") {
			return FormatExport
		}
		if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
			return FormatDotEnv
		}
	}
	return FormatDotEnv
}

// parseDotEnv handles KEY=VALUE style lines.
func parseDotEnv(src string) (*Result, error) {
	vars := make(map[string]string)
	var skipped []string
	for _, line := range strings.Split(src, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			skipped = append(skipped, line)
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.Trim(strings.TrimSpace(line[idx+1:]), `"`)
		if key == "" {
			skipped = append(skipped, line)
			continue
		}
		vars[key] = val
	}
	return &Result{Vars: vars, Skipped: skipped, Format: FormatDotEnv}, nil
}

// parseExport handles "export KEY=VALUE" style lines.
func parseExport(src string) (*Result, error) {
	modified := strings.ReplaceAll(src, "export ", "")
	res, err := parseDotEnv(modified)
	if err != nil {
		return nil, err
	}
	res.Format = FormatExport
	return res, nil
}
