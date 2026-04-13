// Package run provides utilities for executing commands with injected
// environment variable profiles.
package run

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Options configures how a command is run.
type Options struct {
	// Vars are the environment variables to inject.
	Vars map[string]string
	// Inherit controls whether the current process environment is inherited.
	Inherit bool
	// Dir is the working directory for the command. Defaults to current dir.
	Dir string
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Inherit: true,
	}
}

// Run executes the given command with the provided options.
// args[0] is the executable; the rest are arguments.
func Run(args []string, opts Options) error {
	if len(args) == 0 {
		return fmt.Errorf("run: no command specified")
	}

	cmd := exec.Command(args[0], args[1:]...)

	var env []string
	if opts.Inherit {
		env = os.Environ()
	}

	for k, v := range opts.Vars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// BuildEnv merges base OS environment with the provided vars map.
// Profile vars take precedence over inherited values.
func BuildEnv(vars map[string]string, inherit bool) []string {
	result := map[string]string{}

	if inherit {
		for _, entry := range os.Environ() {
			parts := strings.SplitN(entry, "=", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			}
		}
	}

	for k, v := range vars {
		result[k] = v
	}

	out := make([]string, 0, len(result))
	for k, v := range result {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return out
}
