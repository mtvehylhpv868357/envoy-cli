package resolve_test

import (
	"os"
	"testing"

	"github.com/yourorg/envoy-cli/internal/resolve"
)

func TestValue_BasicSubstitution(t *testing.T) {
	vars := map[string]string{"HOST": "localhost", "PORT": "5432"}
	opts := resolve.DefaultOptions()
	out, err := resolve.Value("postgres://${HOST}:${PORT}/db", vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "postgres://localhost:5432/db" {
		t.Errorf("got %q", out)
	}
}

func TestValue_DollarSyntax(t *testing.T) {
	vars := map[string]string{"NAME": "world"}
	opts := resolve.DefaultOptions()
	out, err := resolve.Value("hello $NAME", vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "hello world" {
		t.Errorf("got %q", out)
	}
}

func TestValue_FallbackToOS(t *testing.T) {
	os.Setenv("ENVOY_TEST_OS_VAR", "fromOS")
	defer os.Unsetenv("ENVOY_TEST_OS_VAR")

	vars := map[string]string{}
	opts := resolve.Options{FallbackToOS: true, Strict: false}
	out, err := resolve.Value("${ENVOY_TEST_OS_VAR}", vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "fromOS" {
		t.Errorf("got %q", out)
	}
}

func TestValue_Strict_Unresolved(t *testing.T) {
	vars := map[string]string{}
	opts := resolve.Options{FallbackToOS: false, Strict: true}
	_, err := resolve.Value("${MISSING}", vars, opts)
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestValue_NoRefs_Unchanged(t *testing.T) {
	vars := map[string]string{"X": "y"}
	opts := resolve.DefaultOptions()
	out, err := resolve.Value("plain string", vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "plain string" {
		t.Errorf("got %q", out)
	}
}

func TestVars_ExpandsAll(t *testing.T) {
	vars := map[string]string{
		"BASE_URL": "http://localhost",
		"API_URL":  "${BASE_URL}/api",
	}
	opts := resolve.Options{FallbackToOS: false, Strict: false}
	out, err := resolve.Vars(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_URL"] != "http://localhost/api" {
		t.Errorf("got %q", out["API_URL"])
	}
	if out["BASE_URL"] != "http://localhost" {
		t.Errorf("got %q", out["BASE_URL"])
	}
}

func TestVars_Strict_Error(t *testing.T) {
	vars := map[string]string{"A": "${UNDEFINED}"}
	opts := resolve.Options{FallbackToOS: false, Strict: true}
	_, err := resolve.Vars(vars, opts)
	if err == nil {
		t.Fatal("expected error for unresolved variable in Vars")
	}
}
