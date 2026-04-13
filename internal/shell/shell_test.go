package shell

import (
	"os"
	"strings"
	"testing"
)

func TestDetect_Bash(t *testing.T) {
	os.Setenv("SHELL", "/bin/bash")
	if got := Detect(); got != Bash {
		t.Errorf("expected bash, got %s", got)
	}
}

func TestDetect_Zsh(t *testing.T) {
	os.Setenv("SHELL", "/usr/bin/zsh")
	if got := Detect(); got != Zsh {
		t.Errorf("expected zsh, got %s", got)
	}
}

func TestDetect_Fish(t *testing.T) {
	os.Setenv("SHELL", "/usr/local/bin/fish")
	if got := Detect(); got != Fish {
		t.Errorf("expected fish, got %s", got)
	}
}

func TestExportScript_Bash(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	script := ExportScript(vars, Bash)
	if !strings.Contains(script, "export FOO=") {
		t.Errorf("expected export statement for FOO, got:\n%s", script)
	}
	if !strings.Contains(script, "export BAZ=") {
		t.Errorf("expected export statement for BAZ, got:\n%s", script)
	}
}

func TestExportScript_Fish(t *testing.T) {
	vars := map[string]string{"FOO": "bar"}
	script := ExportScript(vars, Fish)
	if !strings.Contains(script, "set -x FOO") {
		t.Errorf("expected fish set -x, got:\n%s", script)
	}
}

func TestUnsetScript_Bash(t *testing.T) {
	keys := []string{"FOO", "BAR"}
	script := UnsetScript(keys, Bash)
	if !strings.Contains(script, "unset FOO") || !strings.Contains(script, "unset BAR") {
		t.Errorf("expected unset statements, got:\n%s", script)
	}
}

func TestUnsetScript_Fish(t *testing.T) {
	keys := []string{"FOO"}
	script := UnsetScript(keys, Fish)
	if !strings.Contains(script, "set -e FOO") {
		t.Errorf("expected fish set -e, got:\n%s", script)
	}
}
