package sanitize_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/sanitize"
)

func TestDefaultOptions(t *testing.T) {
	opts := sanitize.DefaultOptions()
	if !opts.TrimSpace {
		t.Error("expected TrimSpace to be true by default")
	}
	if !opts.RemoveInvalidKeys {
		t.Error("expected RemoveInvalidKeys to be true by default")
	}
	if opts.RemoveEmpty {
		t.Error("expected RemoveEmpty to be false by default")
	}
}

func TestMap_TrimSpace(t *testing.T) {
	input := map[string]string{"KEY": "  hello world  "}
	opts := sanitize.DefaultOptions()
	out := sanitize.Map(input, opts)
	if got := out["KEY"]; got != "hello world" {
		t.Errorf("expected trimmed value, got %q", got)
	}
}

func TestMap_UppercaseKeys(t *testing.T) {
	input := map[string]string{"my_key": "value"}
	opts := sanitize.DefaultOptions()
	opts.UppercaseKeys = true
	out := sanitize.Map(input, opts)
	if _, ok := out["MY_KEY"]; !ok {
		t.Error("expected key to be uppercased")
	}
	if _, ok := out["my_key"]; ok {
		t.Error("original lowercase key should not be present")
	}
}

func TestMap_RemoveEmpty(t *testing.T) {
	input := map[string]string{"A": "", "B": "val"}
	opts := sanitize.DefaultOptions()
	opts.RemoveEmpty = true
	out := sanitize.Map(input, opts)
	if _, ok := out["A"]; ok {
		t.Error("empty-value key should have been removed")
	}
	if _, ok := out["B"]; !ok {
		t.Error("non-empty key should be retained")
	}
}

func TestMap_RemoveInvalidKeys(t *testing.T) {
	input := map[string]string{
		"1INVALID":  "v",
		"bad-key":   "v",
		"good_KEY":  "v",
		"_PRIVATE":  "v",
	}
	opts := sanitize.DefaultOptions()
	out := sanitize.Map(input, opts)
	if _, ok := out["1INVALID"]; ok {
		t.Error("key starting with digit should be removed")
	}
	if _, ok := out["bad-key"]; ok {
		t.Error("key with hyphen should be removed")
	}
	if _, ok := out["good_KEY"]; !ok {
		t.Error("valid key should be retained")
	}
	if _, ok := out["_PRIVATE"]; !ok {
		t.Error("key starting with underscore should be retained")
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"KEY": "  spaced  "}
	opts := sanitize.DefaultOptions()
	sanitize.Map(input, opts)
	if input["KEY"] != "  spaced  " {
		t.Error("original map should not be mutated")
	}
}

func TestMap_EmptyInput(t *testing.T) {
	out := sanitize.Map(map[string]string{}, sanitize.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
