package validate_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/validate"
)

func TestValidate_RequiredMissing(t *testing.T) {
	schema := validate.Schema{
		Rules: []validate.Rule{
			{Key: "DATABASE_URL", Required: true},
		},
	}
	issues := validate.Validate(map[string]string{}, schema)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Level != "error" {
		t.Errorf("expected error level, got %s", issues[0].Level)
	}
}

func TestValidate_RequiredPresent(t *testing.T) {
	schema := validate.Schema{
		Rules: []validate.Rule{
			{Key: "DATABASE_URL", Required: true},
		},
	}
	issues := validate.Validate(map[string]string{"DATABASE_URL": "postgres://localhost"}, schema)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	schema := validate.Schema{
		Rules: []validate.Rule{
			{Key: "PORT", Pattern: `^\d+$`},
		},
	}
	issues := validate.Validate(map[string]string{"PORT": "abc"}, schema)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Level != "error" {
		t.Errorf("expected error level")
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	schema := validate.Schema{
		Rules: []validate.Rule{
			{Key: "PORT", Pattern: `^\d+$`},
		},
	}
	issues := validate.Validate(map[string]string{"PORT": "8080"}, schema)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestValidate_MinLen(t *testing.T) {
	schema := validate.Schema{
		Rules: []validate.Rule{
			{Key: "SECRET", MinLen: 16},
		},
	}
	issues := validate.Validate(map[string]string{"SECRET": "short"}, schema)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
}

func TestValidate_MaxLen(t *testing.T) {
	schema := validate.Schema{
		Rules: []validate.Rule{
			{Key: "LABEL", MaxLen: 5},
		},
	}
	issues := validate.Validate(map[string]string{"LABEL": "toolongvalue"}, schema)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Level != "warn" {
		t.Errorf("expected warn level, got %s", issues[0].Level)
	}
}

func TestHasErrors_True(t *testing.T) {
	issues := []validate.Issue{{Key: "X", Level: "error", Message: "bad"}}
	if !validate.HasErrors(issues) {
		t.Error("expected HasErrors to be true")
	}
}

func TestHasErrors_False(t *testing.T) {
	issues := []validate.Issue{{Key: "X", Level: "warn", Message: "ok"}}
	if validate.HasErrors(issues) {
		t.Error("expected HasErrors to be false")
	}
}
