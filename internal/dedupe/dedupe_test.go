package dedupe_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/dedupe"
)

func TestMap_KeepLast(t *testing.T) {
	pairs := []string{"FOO=first", "BAR=one", "FOO=second"}
	result := dedupe.Map(pairs, dedupe.DefaultOptions())
	if result["FOO"] != "second" {
		t.Fatalf("expected 'second', got %q", result["FOO"])
	}
}

func TestMap_KeepFirst(t *testing.T) {
	pairs := []string{"FOO=first", "BAR=one", "FOO=second"}
	opts := dedupe.Options{KeepFirst: true}
	result := dedupe.Map(pairs, opts)
	if result["FOO"] != "first" {
		t.Fatalf("expected 'first', got %q", result["FOO"])
	}
}

func TestMap_NoDuplicates(t *testing.T) {
	pairs := []string{"A=1", "B=2", "C=3"}
	result := dedupe.Map(pairs, dedupe.DefaultOptions())
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestMap_SkipsEmptyKey(t *testing.T) {
	pairs := []string{"=value", "VALID=yes"}
	result := dedupe.Map(pairs, dedupe.DefaultOptions())
	if _, ok := result[""]; ok {
		t.Fatal("empty key should be skipped")
	}
	if result["VALID"] != "yes" {
		t.Fatal("expected VALID=yes")
	}
}

func TestDuplicates_ReturnsDupKeys(t *testing.T) {
	pairs := []string{"FOO=1", "BAR=2", "FOO=3", "BAZ=4", "BAR=5"}
	dups := dedupe.Duplicates(pairs)
	if len(dups) != 2 {
		t.Fatalf("expected 2 duplicates, got %d: %v", len(dups), dups)
	}
}

func TestDuplicates_NoneReturnsNil(t *testing.T) {
	pairs := []string{"A=1", "B=2"}
	dups := dedupe.Duplicates(pairs)
	if len(dups) != 0 {
		t.Fatalf("expected no duplicates, got %v", dups)
	}
}

func TestMap_EmptyInput(t *testing.T) {
	result := dedupe.Map(nil, dedupe.DefaultOptions())
	if len(result) != 0 {
		t.Fatal("expected empty map for nil input")
	}
}
