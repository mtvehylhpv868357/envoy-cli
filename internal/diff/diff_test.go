package diff

import (
	"testing"
)

func TestCompare_Added(t *testing.T) {
	before := map[string]string{"FOO": "bar"}
	after := map[string]string{"FOO": "bar", "BAZ": "qux"}

	changes := Compare(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Action != "added" || changes[0].Key != "BAZ" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompare_Removed(t *testing.T) {
	before := map[string]string{"FOO": "bar", "BAZ": "qux"}
	after := map[string]string{"FOO": "bar"}

	changes := Compare(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Action != "removed" || changes[0].Key != "BAZ" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompare_Modified(t *testing.T) {
	before := map[string]string{"FOO": "bar"}
	after := map[string]string{"FOO": "baz"}

	changes := Compare(before, after)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Action != "modified" || c.Old != "bar" || c.New != "baz" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	changes := Compare(env, env)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestCompare_Sorted(t *testing.T) {
	before := map[string]string{}
	after := map[string]string{"Z_VAR": "1", "A_VAR": "2", "M_VAR": "3"}

	changes := Compare(before, after)
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(changes))
	}
	if changes[0].Key != "A_VAR" || changes[1].Key != "M_VAR" || changes[2].Key != "Z_VAR" {
		t.Errorf("changes not sorted: %+v", changes)
	}
}

func TestSummary(t *testing.T) {
	changes := []Change{
		{Action: "added"},
		{Action: "added"},
		{Action: "removed"},
		{Action: "modified"},
	}
	s := Summary(changes)
	if s["added"] != 2 || s["removed"] != 1 || s["modified"] != 1 {
		t.Errorf("unexpected summary: %+v", s)
	}
}
