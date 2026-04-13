package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func findHistoryCmd(name string) bool {
	for _, c := range rootCmd.Commands() {
		if c.Name() == name {
			return true
		}
	}
	return false
}

func TestHistoryCmd_Registered(t *testing.T) {
	if !findHistoryCmd("history") {
		t.Error("expected 'history' command to be registered on rootCmd")
	}
}

func TestHistoryCmd_SubcommandsRegistered(t *testing.T) {
	var histCmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "history" {
			histCmd = c
			break
		}
	}
	if histCmd == nil {
		t.Fatal("history command not found")
	}
	found := false
	for _, sub := range histCmd.Commands() {
		if sub.Name() == "clear" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'clear' subcommand under history")
	}
}

func TestHistoryCmd_NoHistory(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"history"})
	_ = rootCmd.Execute()

	if !strings.Contains(buf.String(), "No history") {
		t.Errorf("expected 'No history' message, got: %q", buf.String())
	}
}

func TestHistoryCmd_HelpText(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"history", "--help"})
	_ = rootCmd.Execute()

	if !strings.Contains(buf.String(), "chronological") {
		t.Errorf("expected help to mention 'chronological', got: %q", buf.String())
	}
}
