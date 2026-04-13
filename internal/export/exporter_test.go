package export_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/export"
)

var sampleVars = map[string]string{
	"APP_ENV":  "production",
	"DB_PASS":  `sec"ret`,
	"LOG_LEVEL": "info",
}

func TestDefaultOptions(t *testing.T) {
	opts := export.DefaultOptions()
	if opts.Format != export.FormatDotEnv {
		t.Errorf("expected dotenv format, got %q", opts.Format)
	}
	if !opts.Sorted {
		t.Error("expected Sorted to be true by default")
	}
}

func TestToBytes_DotEnv(t *testing.T) {
	opts := export.DefaultOptions()
	out, err := export.ToBytes(sampleVars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "APP_ENV=\"") {
		t.Errorf("expected quoted APP_ENV in output, got:\n%s", s)
	}
	if !strings.Contains(s, `DB_PASS="sec\"ret"`) {
		t.Errorf("expected escaped quote in DB_PASS, got:\n%s", s)
	}
}

func TestToBytes_DotEnv_Sorted(t *testing.T) {
	opts := export.DefaultOptions()
	out, err := export.ToBytes(sampleVars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if !strings.HasPrefix(lines[0], "APP_ENV") {
		t.Errorf("expected APP_ENV first (sorted), got %q", lines[0])
	}
}

func TestToBytes_ExportFormat(t *testing.T) {
	opts := export.Options{Format: export.FormatExport, Sorted: false, Quote: false}
	out, err := export.ToBytes(map[string]string{"FOO": "bar"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "export FOO=bar") {
		t.Errorf("expected 'export FOO=bar', got %q", string(out))
	}
}

func TestToBytes_JSON(t *testing.T) {
	opts := export.Options{Format: export.FormatJSON}
	out, err := export.ToBytes(map[string]string{"KEY": "value"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]string
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", parsed["KEY"])
	}
}

func TestToBytes_UnknownFormat(t *testing.T) {
	opts := export.Options{Format: "xml"}
	_, err := export.ToBytes(sampleVars, opts)
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}
