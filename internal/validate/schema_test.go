package validate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-cli/internal/validate"
)

func tempSchemaDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "validate-schema-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSaveAndLoadSchema_RoundTrip(t *testing.T) {
	dir := tempSchemaDir(t)
	path := filepath.Join(dir, "schema.json")

	original := validate.Schema{
		Rules: []validate.Rule{
			{Key: "API_KEY", Required: true, MinLen: 8},
			{Key: "PORT", Pattern: `^\d+$`},
		},
	}

	if err := validate.SaveSchema(path, original); err != nil {
		t.Fatalf("SaveSchema: %v", err)
	}

	loaded, err := validate.LoadSchema(path)
	if err != nil {
		t.Fatalf("LoadSchema: %v", err)
	}

	if len(loaded.Rules) != len(original.Rules) {
		t.Fatalf("expected %d rules, got %d", len(original.Rules), len(loaded.Rules))
	}
	if loaded.Rules[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", loaded.Rules[0].Key)
	}
	if loaded.Rules[0].MinLen != 8 {
		t.Errorf("expected MinLen 8, got %d", loaded.Rules[0].MinLen)
	}
}

func TestLoadSchema_NotFound(t *testing.T) {
	_, err := validate.LoadSchema("/nonexistent/path/schema.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDefaultSchemaPath(t *testing.T) {
	path := validate.DefaultSchemaPath("/some/dir")
	expected := "/some/dir/.envoy-schema.json"
	if path != expected {
		t.Errorf("expected %s, got %s", expected, path)
	}
}
