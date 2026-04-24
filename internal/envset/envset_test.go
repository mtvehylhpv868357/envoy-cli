package envset_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/envset"
)

func TestUnion_MergesAllMaps(t *testing.T) {
	opts := envset.DefaultOptions()
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3", "FOO": "overwritten"}
	out := envset.Union(opts, a, b)
	if out["FOO"] != "overwritten" {
		t.Errorf("expected FOO=overwritten, got %s", out["FOO"])
	}
	if out["BAR"] != "2" || out["BAZ"] != "3" {
		t.Errorf("missing expected keys in union: %v", out)
	}
}

func TestUnion_EmptyMaps(t *testing.T) {
	opts := envset.DefaultOptions()
	out := envset.Union(opts)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestIntersect_KeysInAll(t *testing.T) {
	opts := envset.DefaultOptions()
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"FOO": "99", "BAZ": "3"}
	out := envset.Intersect(opts, a, b)
	if len(out) != 1 || out["FOO"] != "1" {
		t.Errorf("expected only FOO in intersection, got %v", out)
	}
}

func TestIntersect_NoCommonKeys(t *testing.T) {
	opts := envset.DefaultOptions()
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"BAR": "2"}
	out := envset.Intersect(opts, a, b)
	if len(out) != 0 {
		t.Errorf("expected empty intersection, got %v", out)
	}
}

func TestDifference_KeysOnlyInBase(t *testing.T) {
	opts := envset.DefaultOptions()
	base := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	other := map[string]string{"BAR": "x"}
	out := envset.Difference(opts, base, other)
	if _, ok := out["BAR"]; ok {
		t.Error("BAR should be excluded")
	}
	if out["FOO"] != "1" || out["BAZ"] != "3" {
		t.Errorf("expected FOO and BAZ in difference, got %v", out)
	}
}

func TestSymmetricDiff_ExcludesShared(t *testing.T) {
	opts := envset.DefaultOptions()
	a := map[string]string{"FOO": "1", "SHARED": "x"}
	b := map[string]string{"BAR": "2", "SHARED": "x"}
	out := envset.SymmetricDiff(opts, a, b)
	if _, ok := out["SHARED"]; ok {
		t.Error("SHARED should not appear in symmetric diff")
	}
	if out["FOO"] != "1" || out["BAR"] != "2" {
		t.Errorf("expected FOO and BAR in symmetric diff, got %v", out)
	}
}

func TestUnion_CaseInsensitive(t *testing.T) {
	opts := envset.Options{CaseSensitive: false}
	a := map[string]string{"foo": "1"}
	b := map[string]string{"FOO": "2"}
	out := envset.Union(opts, a, b)
	if out["FOO"] != "2" {
		t.Errorf("expected FOO=2 with case-insensitive union, got %v", out)
	}
}
