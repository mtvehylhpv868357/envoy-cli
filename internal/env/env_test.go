package env

import (
	"os"
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		input    string
		wantKey  string
		wantVal  string
		wantOK   bool
	}{
		{"KEY=VALUE", "KEY", "VALUE", true},
		{`KEY="VALUE"`, "KEY", "VALUE", true},
		{"# comment", "", "", false},
		{"", "", "", false},
		{"NOEQUALS", "", "", false},
		{"  SPACED  =  value  ", "SPACED", "value", true},
	}
	for _, tt := range tests {
		k, v, ok := ParseLine(tt.input)
		if ok != tt.wantOK || k != tt.wantKey || v != tt.wantVal {
			t.Errorf("ParseLine(%q) = (%q, %q, %v), want (%q, %q, %v)",
				tt.input, k, v, ok, tt.wantKey, tt.wantVal, tt.wantOK)
		}
	}
}

func TestParseDotEnv(t *testing.T) {
	content := "# comment\nFOO=bar\nBAZ=\"qux\"\n\nEMPTY="
	vars := ParseDotEnv(content)
	if vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", vars["FOO"])
	}
	if vars["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", vars["BAZ"])
	}
	if _, ok := vars["EMPTY"]; ok {
		t.Errorf("expected EMPTY to be skipped (no value)") 
	}
}

func TestApply(t *testing.T) {
	vars := Vars{"TEST_ENVOY_KEY": "hello"}
	if err := Apply(vars); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if got := os.Getenv("TEST_ENVOY_KEY"); got != "hello" {
		t.Errorf("expected TEST_ENVOY_KEY=hello, got %q", got)
	}
	os.Unsetenv("TEST_ENVOY_KEY")
}

func TestExport(t *testing.T) {
	vars := Vars{"MY_VAR": "my_value"}
	out := Export(vars)
	if !strings.Contains(out, "export MY_VAR=") {
		t.Errorf("Export output missing expected line, got: %q", out)
	}
}
