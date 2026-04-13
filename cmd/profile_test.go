package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestProfileCmd_SubcommandsRegistered(t *testing.T) {
	subcmds := map[string]bool{}
	for _, sub := range profileCmd.Commands() {
		subcmds[sub.Use] = true
	}
	expected := []string{"list", "use <name>", "delete <name>"}
	for _, name := range expected {
		if !subcmds[name] {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}

func TestProfileListCmd_NoProfiles(t *testing.T) {
	t.Setenv("ENVOY_HOME", t.TempDir())

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	profileListCmd.SetOut(buf)

	var outBuf bytes.Buffer
	profileListCmd.SetOut(&outBuf)

	err := profileListCmd.RunE(profileListCmd, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProfileUseCmd_MissingArg(t *testing.T) {
	err := profileUseCmd.Args(profileUseCmd, []string{})
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestProfileDeleteCmd_MissingArg(t *testing.T) {
	err := profileDeleteCmd.Args(profileDeleteCmd, []string{})
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestProfileCmd_HelpText(t *testing.T) {
	help := profileCmd.Short
	if !strings.Contains(help, "profile") && !strings.Contains(strings.ToLower(help), "environment") {
		t.Errorf("expected short description to mention profiles or environment, got: %q", help)
	}
}
