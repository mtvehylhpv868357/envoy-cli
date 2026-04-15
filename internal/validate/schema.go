package validate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const defaultSchemaFile = ".envoy-schema.json"

// LoadSchema reads a Schema from a JSON file.
func LoadSchema(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Schema{}, fmt.Errorf("read schema: %w", err)
	}
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return Schema{}, fmt.Errorf("parse schema: %w", err)
	}
	return s, nil
}

// SaveSchema writes a Schema to a JSON file, creating directories as needed.
func SaveSchema(path string, s Schema) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create dirs: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal schema: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// DefaultSchemaPath returns the default schema file path relative to dir.
func DefaultSchemaPath(dir string) string {
	return filepath.Join(dir, defaultSchemaFile)
}
