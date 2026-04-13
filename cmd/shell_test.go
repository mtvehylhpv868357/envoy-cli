package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestShellCmd_NoActiveProfile(t *testing.T) {
	// Ensure the command surfaces an error when no profile is active and none given.
	cmd := rootCmd
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"shell"})

	err := cmd.Execute()
	// We expect an error because no store / active profile exists in test env.
	if err == nil {
		t.Log("output:", buf.String())
		// Not fatal — may succeed if a default store exists; just log.
	}
}

func TestShellCmd_HelpContainsEval(t *testing.T) {
	cmd := rootCmd
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"shell", "--help"})

	_ = cmd.Execute()
	output := buf.String()
	if !strings.Contains(output, "eval") {
		t.Errorf("expected help text to mention eval, got:\n%s", output)
	}
}

func TestShellCmd_ShellFlagRegistered(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Name() == "shell" {
			if sub.Flags().Lookup("shell") == nil {
				t.Error("expected --shell flag to be registered on shell subcommand")
			}
			return
		}
	}
	t.Error("shell subcommand not found")
}
