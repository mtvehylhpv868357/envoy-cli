package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func findAliasCmd(t *testing.T) interface{ Commands() []interface{} } {
	t.Helper()
	for _, c := range rootCmd.Commands() {
		if c.Name() == "alias" {
			return nil
		}
	}
	return nil
}

func TestAliasCmd_Registered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "alias" {
			return
		}
	}
	t.Error("alias command not registered")
}

func TestAliasCmd_SubcommandsRegistered(t *testing.T) {
	var aliasCmd interface{ Commands() []*cobra_command }
	_ = aliasCmd
	for _, c := range rootCmd.Commands() {
		if c.Name() != "alias" {
			continue
		}
		names := map[string]bool{}
		for _, sub := range c.Commands() {
			names[sub.Name()] = true
		}
		for _, want := range []string{"set", "list", "remove"} {
			if !names[want] {
				t.Errorf("missing subcommand: %s", want)
			}
		}
		return
	}
	t.Error("alias command not found")
}

func TestAliasCmd_SetRequiresTwoArgs(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"alias", "set", "onlyone"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing second arg")
	}
}

func TestAliasCmd_HelpContainsAlias(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"alias", "--help"})
	_ = rootCmd.Execute()
	if !strings.Contains(buf.String(), "alias") {
		t.Error("help output does not mention alias")
	}
}

type cobra_command = struct{}
