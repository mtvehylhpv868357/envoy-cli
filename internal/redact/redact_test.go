package redact_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/redact"
)

func TestIsSensitive_MatchesKnownKeys(t *testing.T) {
	sensitive := []string{
		"API_KEY", "DB_PASSWORD", "AUTH_TOKEN", "PRIVATE_KEY",
		"SECRET", "AWS_SECRET_ACCESS_KEY", "GITHUB_TOKEN",
	}
	for _, k := range sensitive {
		if !redact.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	safe := []string{"APP_ENV", "PORT", "HOST", "DEBUG", "LOG_LEVEL"}
	for _, k := range safe {
		if redact.IsSensitive(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}

func TestValue_FullRedact(t *testing.T) {
	opts := redact.DefaultOptions()
	got := redact.Value("supersecretvalue", opts)
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestValue_PartialReveal(t *testing.T) {
	opts := redact.Options{
		Replacement:   "[REDACTED]",
		PartialReveal: true,
		RevealChars:   4,
	}
	got := redact.Value("supersecretvalue", opts)
	expected := "[REDACTED]...alue"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestValue_ShortValue_FullRedact(t *testing.T) {
	opts := redact.Options{
		Replacement:   "[REDACTED]",
		PartialReveal: true,
		RevealChars:   4,
	}
	got := redact.Value("abc", opts)
	if got != "[REDACTED]" {
		t.Errorf("short value should be fully redacted, got %q", got)
	}
}

func TestMap_RedactsSensitiveKeys(t *testing.T) {
	input := map[string]string{
		"API_KEY":  "my-secret-key",
		"APP_ENV":  "production",
		"DB_PASS":  "hunter2",
		"LOG_LEVEL": "info",
	}
	opts := redact.DefaultOptions()
	out := redact.Map(input, opts)

	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be redacted, got %q", out["API_KEY"])
	}
	if out["DB_PASS"] != "[REDACTED]" {
		t.Errorf("DB_PASS should be redacted, got %q", out["DB_PASS"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged, got %q", out["APP_ENV"])
	}
	if out["LOG_LEVEL"] != "info" {
		t.Errorf("LOG_LEVEL should be unchanged, got %q", out["LOG_LEVEL"])
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"SECRET_TOKEN": "abc123"}
	opts := redact.DefaultOptions()
	redact.Map(input, opts)
	if input["SECRET_TOKEN"] != "abc123" {
		t.Error("Map should not mutate the input map")
	}
}
