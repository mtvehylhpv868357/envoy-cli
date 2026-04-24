package envcount_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/envcount"
)

func TestCount_EmptyMap(t *testing.T) {
	r := envcount.Count("dev", map[string]string{})
	if r.Total != 0 || r.Empty != 0 || r.NonEmpty != 0 {
		t.Fatalf("expected all zeros, got %+v", r)
	}
	if r.Profile != "dev" {
		t.Fatalf("expected profile 'dev', got %q", r.Profile)
	}
}

func TestCount_MixedValues(t *testing.T) {
	vars := map[string]string{
		"HOST":  "localhost",
		"PORT":  "8080",
		"DEBUG": "",
	}
	r := envcount.Count("staging", vars)
	if r.Total != 3 {
		t.Fatalf("expected total 3, got %d", r.Total)
	}
	if r.NonEmpty != 2 {
		t.Fatalf("expected non-empty 2, got %d", r.NonEmpty)
	}
	if r.Empty != 1 {
		t.Fatalf("expected empty 1, got %d", r.Empty)
	}
}

func TestCount_AllEmpty(t *testing.T) {
	vars := map[string]string{"A": "", "B": ""}
	r := envcount.Count("prod", vars)
	if r.Empty != 2 || r.NonEmpty != 0 {
		t.Fatalf("unexpected counts: %+v", r)
	}
}

func TestCountMany_SortedByName(t *testing.T) {
	profiles := map[string]map[string]string{
		"prod":    {"A": "1"},
		"dev":     {"B": "2", "C": ""},
		"staging": {},
	}
	results := envcount.CountMany(profiles)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Profile != "dev" || results[1].Profile != "prod" || results[2].Profile != "staging" {
		t.Fatalf("unexpected order: %v", results)
	}
	if results[0].Total != 2 {
		t.Fatalf("dev should have 2 vars, got %d", results[0].Total)
	}
}

func TestTotal_SumsAllCounts(t *testing.T) {
	results := []envcount.Result{
		{Total: 3},
		{Total: 5},
		{Total: 2},
	}
	if got := envcount.Total(results); got != 10 {
		t.Fatalf("expected total 10, got %d", got)
	}
}

func TestResult_String(t *testing.T) {
	r := envcount.Result{Profile: "dev", Total: 4, NonEmpty: 3, Empty: 1}
	s := r.String()
	if s == "" {
		t.Fatal("expected non-empty string from Result.String()")
	}
}
