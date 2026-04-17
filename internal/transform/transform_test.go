package transform

import (
	"testing"
)

func TestMap_Uppercase(t *testing.T) {
	vars := map[string]string{"KEY": "hello", "OTHER": "world"}
	result, err := Map(vars, Options{Op: OpUppercase})
	if err != nil {
		t.Fatal(err)
	}
	if result["KEY"] != "HELLO" || result["OTHER"] != "WORLD" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestMap_Lowercase(t *testing.T) {
	vars := map[string]string{"KEY": "HELLO"}
	result, err := Map(vars, Options{Op: OpLowercase})
	if err != nil {
		t.Fatal(err)
	}
	if result["KEY"] != "hello" {
		t.Errorf("got %q", result["KEY"])
	}
}

func TestMap_TrimSpace(t *testing.T) {
	vars := map[string]string{"KEY": "  spaced  "}
	result, err := Map(vars, Options{Op: OpTrimSpace})
	if err != nil {
		t.Fatal(err)
	}
	if result["KEY"] != "spaced" {
		t.Errorf("got %q", result["KEY"])
	}
}

func TestMap_Base64RoundTrip(t *testing.T) {
	vars := map[string]string{"SECRET": "my-secret-value"}
	encoded, err := Map(vars, Options{Op: OpBase64Encode})
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := Map(encoded, Options{Op: OpBase64Decode})
	if err != nil {
		t.Fatal(err)
	}
	if decoded["SECRET"] != "my-secret-value" {
		t.Errorf("round-trip failed: got %q", decoded["SECRET"])
	}
}

func TestMap_TargetedKeys(t *testing.T) {
	vars := map[string]string{"A": "hello", "B": "world"}
	result, err := Map(vars, Options{Keys: []string{"A"}, Op: OpUppercase})
	if err != nil {
		t.Fatal(err)
	}
	if result["A"] != "HELLO" {
		t.Errorf("expected A to be uppercased")
	}
	if result["B"] != "world" {
		t.Errorf("expected B to be unchanged")
	}
}

func TestMap_UnknownOp_Passthrough(t *testing.T) {
	vars := map[string]string{"X": "value"}
	result, err := Map(vars, Options{Op: Op("unknown")})
	if err != nil {
		t.Fatal(err)
	}
	if result["X"] != "value" {
		t.Errorf("expected passthrough, got %q", result["X"])
	}
}
