package format

import (
	"strings"
	"testing"
)

func TestTable_OutputContainsKeys(t *testing.T) {
	vars := map[string]string{
		"APP_ENV":  "production",
		"APP_PORT": "8080",
	}
	var sb strings.Builder
	Table(&sb, vars, DefaultOptions())
	out := sb.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV in output")
	}
	if !strings.Contains(out, "APP_PORT") {
		t.Error("expected APP_PORT in output")
	}
}

func TestTable_MasksSecretValues(t *testing.T) {
	vars := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_ENV":     "staging",
	}
	var sb strings.Builder
	Table(&sb, vars, DefaultOptions())
	out := sb.String()
	if strings.Contains(out, "supersecret") {
		t.Error("secret value should be masked")
	}
	if !strings.Contains(out, "********") {
		t.Error("expected masked placeholder in output")
	}
	if !strings.Contains(out, "staging") {
		t.Error("non-secret value should be visible")
	}
}

func TestTable_MaskDisabled(t *testing.T) {
	vars := map[string]string{
		"API_TOKEN": "tok_abc123",
	}
	opts := DefaultOptions()
	opts.MaskSecrets = false
	var sb strings.Builder
	Table(&sb, vars, opts)
	out := sb.String()
	if !strings.Contains(out, "tok_abc123") {
		t.Error("value should be visible when masking is disabled")
	}
}

func TestTable_SortedOutput(t *testing.T) {
	vars := map[string]string{
		"ZEBRA": "z",
		"ALPHA": "a",
		"MANGO": "m",
	}
	var sb strings.Builder
	Table(&sb, vars, DefaultOptions())
	out := sb.String()
	alphaIdx := strings.Index(out, "ALPHA")
	mangoIdx := strings.Index(out, "MANGO")
	zebraIdx := strings.Index(out, "ZEBRA")
	if !(alphaIdx < mangoIdx && mangoIdx < zebraIdx) {
		t.Errorf("expected sorted output: ALPHA < MANGO < ZEBRA, got positions %d %d %d", alphaIdx, mangoIdx, zebraIdx)
	}
}

func TestTable_EmptyMap(t *testing.T) {
	var sb strings.Builder
	Table(&sb, map[string]string{}, DefaultOptions())
	out := sb.String()
	if !strings.Contains(out, "KEY") {
		t.Error("header should still be printed for empty map")
	}
}

func TestDefaultOptions_HasSecretPatterns(t *testing.T) {
	opts := DefaultOptions()
	if len(opts.SecretPatterns) == 0 {
		t.Error("expected default secret patterns to be non-empty")
	}
	if !opts.MaskSecrets {
		t.Error("expected MaskSecrets to be true by default")
	}
}
