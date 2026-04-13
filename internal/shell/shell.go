package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ShellType represents a supported shell.
type ShellType string

const (
	Bash ShellType = "bash"
	Zsh  ShellType = "zsh"
	Fish ShellType = "fish"
)

// Detect attempts to detect the current shell from the SHELL environment variable.
func Detect() ShellType {
	shell := os.Getenv("SHELL")
	base := strings.ToLower(filepath.Base(shell))
	switch base {
	case "zsh":
		return Zsh
	case "fish":
		return Fish
	default:
		return Bash
	}
}

// ExportScript returns a shell-specific script that exports the given env vars.
func ExportScript(vars map[string]string, shell ShellType) string {
	var sb strings.Builder
	switch shell {
	case Fish:
		for k, v := range vars {
			fmt.Fprintf(&sb, "set -x %s %q;\n", k, v)
		}
	default:
		// bash / zsh
		for k, v := range vars {
			fmt.Fprintf(&sb, "export %s=%q\n", k, v)
		}
	}
	return sb.String()
}

// UnsetScript returns a shell-specific script that unsets the given variable names.
func UnsetScript(keys []string, shell ShellType) string {
	var sb strings.Builder
	switch shell {
	case Fish:
		for _, k := range keys {
			fmt.Fprintf(&sb, "set -e %s;\n", k)
		}
	default:
		for _, k := range keys {
			fmt.Fprintf(&sb, "unset %s\n", k)
		}
	}
	return sb.String()
}
