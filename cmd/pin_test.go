package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findPinCmd(root *cobra.Command, name string) *cobra.Command {
	for _, sub := range root.Commands() {
		if sub.Name() == "pin" {
			for _, child := range sub.Commands() {
				if child.Name() == name {
					return child
				}
			}
		}
	}
	return nil
}

func TestPinCmd_Registered(t *testing.T) {
	var found bool
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "pin" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'pin' command to be registered")
	}
}

func TestPinCmd_SubcommandsRegistered(t *testing.T) {
	expected := []string{"set", "get", "remove", "list"}
	var pinCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "pin" {
			pinCmd = cmd
			break
		}
	}
	if pinCmd == nil {
		t.Fatal("pin command not found")
	}
	names := make(map[string]bool)
	for _, sub := range pinCmd.Commands() {
		names[sub.Name()] = true
	}
	for _, e := range expected {
		if !names[e] {
			t.Errorf("expected sub-command %q to be registered", e)
		}
	}
}

func TestPinCmd_SetRequiresArg(t *testing.T) {
	setCmd := findPinCmd(rootCmd, "set")
	if setCmd == nil {
		t.Fatal("pin set command not found")
	}
	buf := &bytes.Buffer{}
	setCmd.SetOut(buf)
	setCmd.SetErr(buf)
	err := setCmd.Args(setCmd, []string{})
	if err == nil {
		t.Error("expected error when no args provided to pin set")
	}
}

func TestPinCmd_HelpContainsDirectory(t *testing.T) {
	var pinCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "pin" {
			pinCmd = cmd
			break
		}
	}
	if pinCmd == nil {
		t.Fatal("pin command not found")
	}
	help := pinCmd.Long + pinCmd.Short
	if !strings.Contains(strings.ToLower(help), "director") &&
		!strings.Contains(strings.ToLower(help), "profile") {
		t.Errorf("expected help text to mention directory or profile, got: %q", help)
	}
}
