package run

import (
	"os"
	"strings"
	"testing"
)

func TestRun_NoArgs_ReturnsError(t *testing.T) {
	err := Run([]string{}, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty args, got nil")
	}
}

func TestRun_SimpleCommand(t *testing.T) {
	opts := DefaultOptions()
	opts.Vars = map[string]string{"ENVOY_TEST_VAR": "hello"}

	err := Run([]string{"env"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_InvalidCommand(t *testing.T) {
	err := Run([]string{"__nonexistent_command_xyz__"}, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
}

func TestBuildEnv_InheritTrue(t *testing.T) {
	os.Setenv("ENVOY_INHERIT_TEST", "base")
	defer os.Unsetenv("ENVOY_INHERIT_TEST")

	vars := map[string]string{"MY_VAR": "injected"}
	env := BuildEnv(vars, true)

	found := map[string]bool{}
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			found[parts[0]] = true
		}
	}

	if !found["ENVOY_INHERIT_TEST"] {
		t.Error("expected inherited env var ENVOY_INHERIT_TEST to be present")
	}
	if !found["MY_VAR"] {
		t.Error("expected injected var MY_VAR to be present")
	}
}

func TestBuildEnv_InheritFalse(t *testing.T) {
	os.Setenv("ENVOY_INHERIT_TEST", "base")
	defer os.Unsetenv("ENVOY_INHERIT_TEST")

	vars := map[string]string{"MY_VAR": "only"}
	env := BuildEnv(vars, false)

	for _, e := range env {
		if strings.HasPrefix(e, "ENVOY_INHERIT_TEST=") {
			t.Error("inherited var should not appear when inherit=false")
		}
	}

	found := false
	for _, e := range env {
		if e == "MY_VAR=only" {
			found = true
		}
	}
	if !found {
		t.Error("expected MY_VAR=only in result")
	}
}

func TestBuildEnv_ProfileOverridesInherited(t *testing.T) {
	os.Setenv("OVERRIDE_ME", "original")
	defer os.Unsetenv("OVERRIDE_ME")

	vars := map[string]string{"OVERRIDE_ME": "new_value"}
	env := BuildEnv(vars, true)

	for _, e := range env {
		if e == "OVERRIDE_ME=original" {
			t.Error("profile var should override inherited value")
		}
	}
}
