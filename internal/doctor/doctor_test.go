package doctor_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/doctor"
)

func TestCheck_EmptyValue_Warning(t *testing.T) {
	vars := map[string]string{
		"MY_VAR": "",
	}
	report := doctor.Check(vars)
	if len(report.Findings) == 0 {
		t.Fatal("expected at least one finding for empty value")
	}
	found := false
	for _, f := range report.Findings {
		if f.Key == "MY_VAR" && f.Severity == doctor.SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected warning finding for MY_VAR empty value")
	}
}

func TestCheck_LowercaseKey_Warning(t *testing.T) {
	vars := map[string]string{
		"my_var": "somevalue",
	}
	report := doctor.Check(vars)
	found := false
	for _, f := range report.Findings {
		if f.Key == "my_var" && f.Severity == doctor.SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for lowercase key")
	}
}

func TestCheck_ShortSecret_Error(t *testing.T) {
	vars := map[string]string{
		"API_KEY": "abc",
	}
	report := doctor.Check(vars)
	if !report.HasErrors() {
		t.Error("expected error finding for short secret value")
	}
}

func TestCheck_ValidVars_NoErrors(t *testing.T) {
	vars := map[string]string{
		"DATABASE_URL": "postgres://localhost:5432/mydb",
		"APP_ENV":      "production",
	}
	report := doctor.Check(vars)
	if report.HasErrors() {
		t.Errorf("expected no errors, got: %v", report.Findings)
	}
}

func TestReport_Summary(t *testing.T) {
	vars := map[string]string{
		"API_KEY": "x",
		"low_key": "value",
	}
	report := doctor.Check(vars)
	summary := report.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestCheck_UnsetOSRef_Warning(t *testing.T) {
	vars := map[string]string{
		"MY_TOKEN": "$DEFINITELY_NOT_SET_XYZ_123",
	}
	report := doctor.Check(vars)
	found := false
	for _, f := range report.Findings {
		if f.Key == "MY_TOKEN" && f.Severity == doctor.SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for unset OS variable reference")
	}
}
