package interpolate

import (
	"os"
	"testing"
)

func osLookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func TestValue_BasicSubstitution(t *testing.T) {
	vars := map[string]string{"HOST": "localhost"}
	out, err := Value("http://${HOST}:8080", vars, nil, DefaultOptions())
	if err != nil || out != "http://localhost:8080" {
		t.Fatalf("expected http://localhost:8080, got %q err=%v", out, err)
	}
}

func TestValue_DollarSyntax(t *testing.T) {
	vars := map[string]string{"PORT": "3000"}
	out, err := Value("$PORT", vars, nil, DefaultOptions())
	if err != nil || out != "3000" {
		t.Fatalf("expected 3000, got %q err=%v", out, err)
	}
}

func TestValue_DefaultFallback(t *testing.T) {
	out, err := Value("${MISSING:-fallback}", nil, nil, DefaultOptions())
	if err != nil || out != "fallback" {
		t.Fatalf("expected fallback, got %q err=%v", out, err)
	}
}

func TestValue_LookupFallback(t *testing.T) {
	t.Setenv("ENVOY_TEST_VAR", "from-env")
	out, err := Value("${ENVOY_TEST_VAR}", nil, osLookup, DefaultOptions())
	if err != nil || out != "from-env" {
		t.Fatalf("expected from-env, got %q err=%v", out, err)
	}
}

func TestValue_Strict_Unresolved(t *testing.T) {
	opts := Options{Strict: true}
	_, err := Value("${UNDEFINED}", nil, nil, opts)
	if err == nil {
		t.Fatal("expected error for unresolved variable in strict mode")
	}
}

func TestValue_NoRefs_Unchanged(t *testing.T) {
	out, err := Value("plain-value", nil, nil, DefaultOptions())
	if err != nil || out != "plain-value" {
		t.Fatalf("expected plain-value, got %q", out)
	}
}

func TestMap_InterpolatesAll(t *testing.T) {
	vars := map[string]string{
		"BASE": "https://example.com",
		"URL":  "${BASE}/api",
	}
	out, err := Map(vars, nil, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if out["URL"] != "https://example.com/api" {
		t.Fatalf("unexpected URL: %s", out["URL"])
	}
}

func TestMap_StrictError(t *testing.T) {
	vars := map[string]string{"X": "${NOPE}"}
	_, err := Map(vars, nil, Options{Strict: true})
	if err == nil {
		t.Fatal("expected error")
	}
}
