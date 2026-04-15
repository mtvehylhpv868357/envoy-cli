package mask_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envoy-cli/internal/mask"
)

func TestIsSensitive_MatchesKnownKeys(t *testing.T) {
	opts := mask.DefaultOptions()
	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "PRIVATE_KEY", "SECRET"}
	for _, k := range sensitive {
		if !opts.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	opts := mask.DefaultOptions()
	safe := []string{"HOST", "PORT", "APP_NAME", "LOG_LEVEL"}
	for _, k := range safe {
		if opts.IsSensitive(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}

func TestValue_MasksWithReveal(t *testing.T) {
	opts := mask.DefaultOptions() // RevealChars = 4
	result := opts.Value("mysupersecret")
	if !strings.HasSuffix(result, "cret") {
		t.Errorf("expected suffix 'cret', got %q", result)
	}
	if !strings.HasPrefix(result, "*") {
		t.Errorf("expected masked prefix, got %q", result)
	}
}

func TestValue_ShortValue_NoReveal(t *testing.T) {
	opts := mask.DefaultOptions()
	opts.RevealChars = 4
	result := opts.Value("abc")
	// len("abc") == 3 < RevealChars(4), so reveal=0
	if result != "***" {
		t.Errorf("expected '***', got %q", result)
	}
}

func TestValue_EmptyString(t *testing.T) {
	opts := mask.DefaultOptions()
	if got := opts.Value(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestVars_MasksSensitiveOnly(t *testing.T) {
	opts := mask.DefaultOptions()
	input := map[string]string{
		"DB_PASSWORD": "supersecret",
		"APP_NAME":    "myapp",
		"API_KEY":     "abcdefgh",
	}
	result := opts.Vars(input)

	if result["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should not be masked, got %q", result["APP_NAME"])
	}
	if result["DB_PASSWORD"] == "supersecret" {
		t.Error("DB_PASSWORD should be masked")
	}
	if result["API_KEY"] == "abcdefgh" {
		t.Error("API_KEY should be masked")
	}
}

func TestVars_DoesNotMutateOriginal(t *testing.T) {
	opts := mask.DefaultOptions()
	input := map[string]string{"SECRET_KEY": "plaintext"}
	_ = opts.Vars(input)
	if input["SECRET_KEY"] != "plaintext" {
		t.Error("original map should not be modified")
	}
}
