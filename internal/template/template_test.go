package template

import (
	"os"
	"testing"
)

func TestRender_BasicSubstitution(t *testing.T) {
	vars := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	res := Render("env=${APP_ENV} port=$PORT", vars)
	if res.Rendered != "env=production port=8080" {
		t.Errorf("unexpected rendered output: %s", res.Rendered)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing vars, got: %v", res.Missing)
	}
}

func TestRender_MissingVars(t *testing.T) {
	vars := map[string]string{"APP_ENV": "staging"}
	res := Render("env=${APP_ENV} db=${DB_HOST}", vars)
	if len(res.Missing) != 1 || res.Missing[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST in missing, got: %v", res.Missing)
	}
	if res.Rendered != "env=staging db=${DB_HOST}" {
		t.Errorf("unexpected rendered output: %s", res.Rendered)
	}
}

func TestRender_DedupMissing(t *testing.T) {
	res := Render("${MISSING} and ${MISSING} again", map[string]string{})
	if len(res.Missing) != 1 {
		t.Errorf("expected 1 unique missing var, got: %d", len(res.Missing))
	}
}

func TestRenderFromEnv(t *testing.T) {
	os.Setenv("ENVOY_TEST_VAR", "hello")
	defer os.Unsetenv("ENVOY_TEST_VAR")
	res := RenderFromEnv("value=${ENVOY_TEST_VAR}")
	if res.Rendered != "value=hello" {
		t.Errorf("unexpected: %s", res.Rendered)
	}
}

func TestValidatePlaceholders_Valid(t *testing.T) {
	if err := ValidatePlaceholders("${APP_ENV} $PORT"); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidatePlaceholders_Invalid(t *testing.T) {
	if err := ValidatePlaceholders("$lowercase"); err == nil {
		t.Error("expected error for lowercase placeholder")
	}
}

func TestRender_EmptyTemplate(t *testing.T) {
	res := Render("", map[string]string{"FOO": "bar"})
	if res.Rendered != "" {
		t.Errorf("expected empty string, got: %s", res.Rendered)
	}
}
