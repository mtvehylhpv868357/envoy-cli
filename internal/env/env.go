package env

import (
	"fmt"
	"os"
	"strings"
)

// Vars represents a map of environment variable key-value pairs.
type Vars map[string]string

// Apply sets all variables in the map as environment variables in the current process.
func Apply(vars Vars) error {
	for k, v := range vars {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("failed to set env var %q: %w", k, err)
		}
	}
	return nil
}

// Export returns a shell-compatible export string for the given vars.
func Export(vars Vars) string {
	var sb strings.Builder
	for k, v := range vars {
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	return sb.String()
}

// ParseLine parses a single "KEY=VALUE" line into a key and value.
// Lines starting with '#' are treated as comments and skipped.
func ParseLine(line string) (key, value string, ok bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key = strings.TrimSpace(parts[0])
	value = strings.Trim(strings.TrimSpace(parts[1]), `"`)
	if key == "" {
		return "", "", false
	}
	return key, value, true
}

// ParseDotEnv parses a .env-style string into a Vars map.
func ParseDotEnv(content string) Vars {
	vars := make(Vars)
	for _, line := range strings.Split(content, "\n") {
		if k, v, ok := ParseLine(line); ok {
			vars[k] = v
		}
	}
	return vars
}
