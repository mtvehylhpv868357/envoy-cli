package lint_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/lint"
)

func TestCheck_ValidVars_NoIssues(t *testing.T) {
	vars := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"PORT":         "8080",
	}
	issues := lint.Check(vars)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestCheck_LowercaseKey_WarnIssue(t *testing.T) {
	vars := map[string]string{
		"database_url": "postgres://localhost/mydb",
	}
	issues := lint.Check(vars)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "warn" {
		t.Errorf("expected warn severity, got %s", issues[0].Severity)
	}
}

func TestCheck_EmptyValue_WarnIssue(t *testing.T) {
	vars := map[string]string{
		"MY_VAR": "",
	}
	issues := lint.Check(vars)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "warn" {
		t.Errorf("expected warn severity, got %s", issues[0].Severity)
	}
}

func TestCheck_SecretKey_ErrorIssue(t *testing.T) {
	vars := map[string]string{
		"API_SECRET": "supersecret",
	}
	issues := lint.Check(vars)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}

func TestHasErrors_WithError(t *testing.T) {
	issues := []lint.Issue{
		{Key: "X", Severity: "warn", Message: "empty value"},
		{Key: "Y", Severity: "error", Message: "secret"},
	}
	if !lint.HasErrors(issues) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_OnlyWarns(t *testing.T) {
	issues := []lint.Issue{
		{Key: "X", Severity: "warn", Message: "empty value"},
	}
	if lint.HasErrors(issues) {
		t.Error("expected HasErrors to return false")
	}
}

func TestIssue_String(t *testing.T) {
	i := lint.Issue{Key: "MY_KEY", Severity: "warn", Message: "some message"}
	s := i.String()
	if s != "[WARN] MY_KEY: some message" {
		t.Errorf("unexpected string: %s", s)
	}
}
