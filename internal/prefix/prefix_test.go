package prefix

import (
	"testing"
)

func TestAdd_BasicPrefix(t *testing.T) {
	vars := map[string]string{"HOST": "localhost", "PORT": "5432"}
	out, errs := Add(vars, "APP_", DefaultOptions())
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", out["APP_HOST"])
	}
	if out["APP_PORT"] != "5432" {
		t.Errorf("expected APP_PORT=5432, got %q", out["APP_PORT"])
	}
}

func TestAdd_EmptyPrefix(t *testing.T) {
	vars := map[string]string{"KEY": "val"}
	out, errs := Add(vars, "", DefaultOptions())
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if out["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", out["KEY"])
	}
}

func TestStrip_RemovesPrefix(t *testing.T) {
	vars := map[string]string{"APP_HOST": "localhost", "APP_PORT": "5432"}
	out := Strip(vars, "APP_")
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", out["PORT"])
	}
}

func TestStrip_NoMatchPassthrough(t *testing.T) {
	vars := map[string]string{"OTHER_KEY": "value"}
	out := Strip(vars, "APP_")
	if out["OTHER_KEY"] != "value" {
		t.Errorf("expected OTHER_KEY=value, got %q", out["OTHER_KEY"])
	}
}

func TestFilterByPrefix_ReturnsMatching(t *testing.T) {
	vars := map[string]string{
		"APP_HOST": "localhost",
		"DB_HOST":  "db",
		"APP_PORT": "8080",
	}
	out := FilterByPrefix(vars, "APP_")
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should not be in filtered result")
	}
}

func TestFilterByPrefix_EmptyResult(t *testing.T) {
	vars := map[string]string{"HOST": "localhost"}
	out := FilterByPrefix(vars, "NOPE_")
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
