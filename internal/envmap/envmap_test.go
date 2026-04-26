package envmap_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/envmap"
)

func TestFromEnviron_Basic(t *testing.T) {
	input := []string{"FOO=bar", "BAZ=qux", "EMPTY="}
	m := envmap.FromEnviron(input)
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", m["FOO"])
	}
	if m["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", m["BAZ"])
	}
	if v, ok := m["EMPTY"]; !ok || v != "" {
		t.Errorf("expected EMPTY=\"\", got %q", v)
	}
}

func TestFromEnviron_SkipsEmptyKey(t *testing.T) {
	input := []string{"=VALUE", "VALID=yes"}
	m := envmap.FromEnviron(input)
	if _, ok := m[""]; ok {
		t.Error("empty key should not be stored")
	}
	if m["VALID"] != "yes" {
		t.Errorf("expected VALID=yes, got %q", m["VALID"])
	}
}

func TestToEnviron_Sorted(t *testing.T) {
	m := map[string]string{"Z": "last", "A": "first", "M": "mid"}
	env := envmap.ToEnviron(m)
	if len(env) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(env))
	}
	if env[0] != "A=first" {
		t.Errorf("expected A=first first, got %q", env[0])
	}
	if env[2] != "Z=last" {
		t.Errorf("expected Z=last last, got %q", env[2])
	}
}

func TestMerge_Overwrite(t *testing.T) {
	dst := map[string]string{"A": "old", "B": "keep"}
	src := map[string]string{"A": "new", "C": "added"}
	opts := envmap.DefaultOptions()
	out := envmap.Merge(dst, src, opts)
	if out["A"] != "new" {
		t.Errorf("expected A=new, got %q", out["A"])
	}
	if out["B"] != "keep" {
		t.Errorf("expected B=keep, got %q", out["B"])
	}
	if out["C"] != "added" {
		t.Errorf("expected C=added, got %q", out["C"])
	}
}

func TestMerge_NoOverwrite(t *testing.T) {
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "should-not-replace"}
	opts := envmap.Options{Overwrite: false}
	out := envmap.Merge(dst, src, opts)
	if out["A"] != "original" {
		t.Errorf("expected A=original, got %q", out["A"])
	}
}

func TestMerge_UppercaseKeys(t *testing.T) {
	dst := map[string]string{}
	src := map[string]string{"foo": "bar", "baz": "qux"}
	opts := envmap.Options{Overwrite: true, UppercaseKeys: true}
	out := envmap.Merge(dst, src, opts)
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
	if _, ok := out["foo"]; ok {
		t.Error("lowercase key should not be present when UppercaseKeys=true")
	}
}

func TestKeys_Sorted(t *testing.T) {
	m := map[string]string{"Z": "1", "A": "2", "M": "3"}
	keys := envmap.Keys(m)
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("unexpected key order: %v", keys)
	}
}
