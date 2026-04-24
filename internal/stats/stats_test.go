package stats

import (
	"testing"
)

func makeProfiles() map[string]map[string]string {
	return map[string]map[string]string{
		"dev": {
			"APP_NAME": "myapp",
			"DB_HOST":  "localhost",
			"DB_PASS":  "",
			"API_KEY":  "abc123",
		},
		"prod": {
			"APP_NAME": "myapp",
			"DB_HOST":  "prod-db",
			"DB_PASS":  "secret",
			"REGION":   "us-east-1",
		},
	}
}

func TestCompute_ProfileCount(t *testing.T) {
	r := Compute(makeProfiles())
	if r.ProfileCount != 2 {
		t.Errorf("expected 2 profiles, got %d", r.ProfileCount)
	}
}

func TestCompute_TotalVars(t *testing.T) {
	r := Compute(makeProfiles())
	if r.TotalVars != 8 {
		t.Errorf("expected 8 total vars, got %d", r.TotalVars)
	}
}

func TestCompute_UniqueAndSharedKeys(t *testing.T) {
	r := Compute(makeProfiles())
	// APP_NAME, DB_HOST, DB_PASS appear in both => 3 shared
	// API_KEY, REGION appear in one each => 2 unique
	if r.SharedKeys != 3 {
		t.Errorf("expected 3 shared keys, got %d", r.SharedKeys)
	}
	if r.UniqueKeys != 2 {
		t.Errorf("expected 2 unique keys, got %d", r.UniqueKeys)
	}
}

func TestCompute_EmptyValues(t *testing.T) {
	r := Compute(makeProfiles())
	// dev.DB_PASS is empty
	if r.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", r.EmptyValues)
	}
}

func TestCompute_SensitiveKeys(t *testing.T) {
	r := Compute(makeProfiles())
	// DB_PASS and API_KEY are sensitive (distinct)
	if r.SensitiveKeys != 2 {
		t.Errorf("expected 2 sensitive keys, got %d", r.SensitiveKeys)
	}
}

func TestCompute_EmptyProfiles(t *testing.T) {
	r := Compute(map[string]map[string]string{})
	if r.ProfileCount != 0 || r.TotalVars != 0 {
		t.Errorf("expected zero stats for empty input")
	}
}

func TestTopKeys_ReturnsTopN(t *testing.T) {
	r := Compute(makeProfiles())
	top := TopKeys(r, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 top keys, got %d", len(top))
	}
	// The top keys should be those appearing in both profiles
	shared := map[string]bool{"APP_NAME": true, "DB_HOST": true, "DB_PASS": true}
	for _, k := range top {
		if !shared[k] {
			t.Errorf("unexpected top key: %s", k)
		}
	}
}

func TestTopKeys_NLargerThanKeys(t *testing.T) {
	r := Compute(makeProfiles())
	top := TopKeys(r, 100)
	if len(top) != 5 {
		t.Errorf("expected 5 keys, got %d", len(top))
	}
}
