package flatten_test

import (
	"testing"

	"github.com/yourusername/envoy-cli/internal/flatten"
)

func TestDefaultOptions(t *testing.T) {
	opts := flatten.DefaultOptions()
	if opts.Delimiter != "_" {
		t.Errorf("expected delimiter '_', got %q", opts.Delimiter)
	}
	if !opts.Uppercase {
		t.Error("expected Uppercase to be true")
	}
}

func TestMap_FlatInput(t *testing.T) {
	input := map[string]any{"key": "value"}
	out, err := flatten.Map(input, flatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", out["KEY"])
	}
}

func TestMap_NestedKeys(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": "5432",
		},
	}
	out, err := flatten.Map(input, flatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
}

func TestMap_WithPrefix(t *testing.T) {
	opts := flatten.DefaultOptions()
	opts.Prefix = "app"
	input := map[string]any{"name": "envoy"}
	out, err := flatten.Map(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "envoy" {
		t.Errorf("expected APP_NAME=envoy, got %q", out["APP_NAME"])
	}
}

func TestMap_NilValue(t *testing.T) {
	input := map[string]any{"empty": nil}
	out, err := flatten.Map(input, flatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["EMPTY"] != "" {
		t.Errorf("expected empty string, got %q", out["EMPTY"])
	}
}

func TestMap_LowercaseOption(t *testing.T) {
	opts := flatten.DefaultOptions()
	opts.Uppercase = false
	input := map[string]any{"Key": "val"}
	out, err := flatten.Map(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["Key"] != "val" {
		t.Errorf("expected Key=val, got %v", out)
	}
}
