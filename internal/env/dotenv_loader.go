package env

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadFromFile reads a .env-style file from disk and returns the parsed Vars.
func LoadFromFile(path string) (Vars, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading env file %q: %w", path, err)
	}
	return ParseDotEnv(string(data)), nil
}

// WriteToFile serializes Vars to a .env-style file at the given path.
// Existing files are overwritten.
func WriteToFile(path string, vars Vars) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating directories for %q: %w", path, err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating env file %q: %w", path, err)
	}
	defer f.Close()
	for k, v := range vars {
		if _, err := fmt.Fprintf(f, "%s=%q\n", k, v); err != nil {
			return fmt.Errorf("writing env var %q: %w", k, err)
		}
	}
	return nil
}

// MergeVars merges src into dst, with src values taking precedence.
func MergeVars(dst, src Vars) Vars {
	result := make(Vars, len(dst)+len(src))
	for k, v := range dst {
		result[k] = v
	}
	for k, v := range src {
		result[k] = v
	}
	return result
}
