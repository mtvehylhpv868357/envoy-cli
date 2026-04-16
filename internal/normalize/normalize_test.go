package normalize_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/normalize"
)

func TestDefaultOptions(t *testing.T) {
	opts := normalize.DefaultOptions()
	if !opts.UppercaseKeys {
		t.Error("expected UppercaseKeys to be true")
	}
	if !opts.TrimSpace {
		t.Error("expected TrimSpace to be true")
	}
	if !opts.ReplaceHyphens {
		t.Error("expected ReplaceHyphens to be true")
	}
}

func TestMap_UppercaseKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	out, changed := normalize.Map(env, normalize.DefaultOptions())
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected key DB_HOST")
	}
	if len(changed) != 1 || changed[0] != "db_host" {
		t.Errorf("unexpected changed list: %v", changed)
	}
}

func TestMap_TrimSpace(t *testing.T) {
	env := map[string]string{"  KEY  ": "  value  "}
	out, changed := normalize.Map(env, normalize.DefaultOptions())
	if v, ok := out["KEY"]; !ok || v != "value" {
		t.Errorf("unexpected result: %v", out)
	}
	if len(changed) == 0 {
		t.Error("expected changed to be non-empty")
	}
}

func TestMap_ReplaceHyphens(t *testing.T) {
	env := map[string]string{"my-key": "val"}
	out, _ := normalize.Map(env, normalize.DefaultOptions())
	if _, ok := out["MY_KEY"]; !ok {
		t.Errorf("expected MY_KEY, got %v", out)
	}
}

func TestMap_NoChanges(t *testing.T) {
	env := map[string]string{"CLEAN": "value"}
	_, changed := normalize.Map(env, normalize.DefaultOptions())
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}

func TestKey_Normalize(t *testing.T) {
	opts := normalize.DefaultOptions()
	result := normalize.Key(" my-key ", opts)
	if result != "MY_KEY" {
		t.Errorf("expected MY_KEY, got %s", result)
	}
}

func TestMap_DisabledOptions(t *testing.T) {
	opts := normalize.Options{UppercaseKeys: false, TrimSpace: false, ReplaceHyphens: false}
	env := map[string]string{"my-key": " val "}
	out, changed := normalize.Map(env, opts)
	if v := out["my-key"]; v != " val " {
		t.Errorf("unexpected value: %q", v)
	}
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}
