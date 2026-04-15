package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func findRenameCmd(t *testing.T) bool {
	t.Helper()
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rename" {
			return true
		}
	}
	return false
}

func TestRenameCmd_Registered(t *testing.T) {
	if !findRenameCmd(t) {
		t.Fatal("rename command not registered on root")
	}
}

func TestRenameCmd_RequiresTwoArgs(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"rename"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestRenameCmd_HelpContainsOldNew(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"rename", "--help"})
	_ = rootCmd.Execute()

	out := buf.String()
	if !strings.Contains(out, "old-name") || !strings.Contains(out, "new-name") {
		t.Errorf("help text should mention old-name and new-name, got:\n%s", out)
	}
}

func TestRenameCmd_OverwriteFlagRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rename" {
			if f := c.Flags().Lookup("overwrite"); f == nil {
				t.Error("expected --overwrite flag on rename command")
			}
			return
		}
	}
	t.Fatal("rename command not found")
}

func TestRenameCmd_StoreFlagRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "rename" {
			if f := c.Flags().Lookup("store"); f == nil {
				t.Error("expected --store flag on rename command")
			}
			return
		}
	}
	t.Fatal("rename command not found")
}
