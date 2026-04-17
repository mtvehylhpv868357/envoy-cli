package filter

import (
	"testing"
)

var sampleVars = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"APP_PORT":    "8080",
	"APP_DEBUG":   "true",
	"LOG_LEVEL":   "info",
}

func TestMap_KeyPattern(t *testing.T) {
	opts := DefaultOptions()
	opts.KeyPattern = `^DB_`
	out, err := Map(sampleVars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestMap_Prefix(t *testing.T) {
	opts := DefaultOptions()
	opts.Prefix = "APP_"
	out, err := Map(sampleVars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestMap_ValuePattern(t *testing.T) {
	opts := DefaultOptions()
	opts.ValuePattern = `^[0-9]+$`
	out, err := Map(sampleVars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out["APP_PORT"] != "8080" {
		t.Error("expected APP_PORT=8080")
	}
}

func TestMap_Invert(t *testing.T) {
	opts := DefaultOptions()
	opts.Prefix = "DB_"
	opts.Invert = true
	out, err := Map(sampleVars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should be excluded when inverted")
	}
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestMap_InvalidKeyPattern(t *testing.T) {
	opts := DefaultOptions()
	opts.KeyPattern = `[invalid`
	_, err := Map(sampleVars, opts)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestMap_NoFilters_ReturnsAll(t *testing.T) {
	out, err := Map(sampleVars, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(sampleVars) {
		t.Fatalf("expected %d, got %d", len(sampleVars), len(out))
	}
}
