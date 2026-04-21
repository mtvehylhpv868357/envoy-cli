package envdiff_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/envdiff"
)

func TestCompare_Missing(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar"}

	res := envdiff.Compare(a, b, false)
	missing := res.Missing()
	if len(missing) != 1 || missing[0].Key != "BAZ" {
		t.Fatalf("expected BAZ missing, got %+v", missing)
	}
}

func TestCompare_Extra(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "EXTRA": "val"}

	res := envdiff.Compare(a, b, true)
	extra := res.Extra()
	if len(extra) != 1 || extra[0].Key != "EXTRA" {
		t.Fatalf("expected EXTRA extra, got %+v", extra)
	}
}

func TestCompare_ExtraIgnoredWhenFalse(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "EXTRA": "val"}

	res := envdiff.Compare(a, b, false)
	if len(res.Extra()) != 0 {
		t.Fatal("expected no extra entries when includeExtra=false")
	}
}

func TestCompare_Changed(t *testing.T) {
	a := map[string]string{"FOO": "original"}
	b := map[string]string{"FOO": "modified"}

	res := envdiff.Compare(a, b, false)
	changed := res.Changed()
	if len(changed) != 1 {
		t.Fatalf("expected 1 changed entry, got %d", len(changed))
	}
	if changed[0].ProfileVal != "original" || changed[0].EnvVal != "modified" {
		t.Fatalf("unexpected values: %+v", changed[0])
	}
}

func TestCompare_Matching(t *testing.T) {
	a := map[string]string{"FOO": "same"}
	b := map[string]string{"FOO": "same"}

	res := envdiff.Compare(a, b, false)
	if len(res.Entries) != 1 || res.Entries[0].Status != envdiff.StatusMatching {
		t.Fatalf("expected matching entry, got %+v", res.Entries)
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	a := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}
	b := map[string]string{}

	res := envdiff.Compare(a, b, false)
	keys := make([]string, len(res.Entries))
	for i, e := range res.Entries {
		keys[i] = e.Key
	}
	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range expected {
		if keys[i] != k {
			t.Fatalf("expected sorted key %s at index %d, got %s", k, i, keys[i])
		}
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	res := envdiff.Compare(map[string]string{}, map[string]string{}, true)
	if len(res.Entries) != 0 {
		t.Fatalf("expected empty result, got %+v", res.Entries)
	}
}
