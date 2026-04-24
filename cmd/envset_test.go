package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func findEnvsetCmd(root *cobra.Command) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Use == "envset" {
			return c
		}
	}
	return nil
}

func TestEnvsetCmd_Registered(t *testing.T) {
	if findEnvsetCmd(rootCmd) == nil {
		t.Fatal("envset command not registered")
	}
}

func TestEnvsetCmd_SubcommandsRegistered(t *testing.T) {
	cmd := findEnvsetCmd(rootCmd)
	if cmd == nil {
		t.Fatal("envset command not found")
	}
	want := []string{"union", "intersect", "diff"}
	found := map[string]bool{}
	for _, sub := range cmd.Commands() {
		for _, w := range want {
			if sub.Name() == w {
				found[w] = true
			}
		}
	}
	for _, w := range want {
		if !found[w] {
			t.Errorf("subcommand %q not registered under envset", w)
		}
	}
}

func TestEnvsetCmd_UnionRequiresTwoArgs(t *testing.T) {
	cmd := findEnvsetCmd(rootCmd)
	if cmd == nil {
		t.Fatal("envset command not found")
	}
	var unionCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "union" {
			unionCmd = sub
		}
	}
	if unionCmd == nil {
		t.Fatal("union subcommand not found")
	}
	err := unionCmd.Args(unionCmd, []string{"only-one"})
	if err == nil {
		t.Error("expected error when fewer than 2 args provided to union")
	}
}

func TestEnvsetCmd_DiffRequiresExactlyTwoArgs(t *testing.T) {
	cmd := findEnvsetCmd(rootCmd)
	if cmd == nil {
		t.Fatal("envset command not found")
	}
	var diffCmd *cobra.Command
	for _, sub := range cmd.Commands() {
		if sub.Name() == "diff" {
			diffCmd = sub
		}
	}
	if diffCmd == nil {
		t.Fatal("diff subcommand not found")
	}
	if err := diffCmd.Args(diffCmd, []string{"a", "b", "c"}); err == nil {
		t.Error("expected error when more than 2 args provided to diff")
	}
	if err := diffCmd.Args(diffCmd, []string{"a", "b"}); err != nil {
		t.Errorf("unexpected error for valid 2-arg diff: %v", err)
	}
}

func TestEnvsetCmd_StoreFlagRegistered(t *testing.T) {
	cmd := findEnvsetCmd(rootCmd)
	if cmd == nil {
		t.Fatal("envset command not found")
	}
	if cmd.PersistentFlags().Lookup("store") == nil {
		t.Error("--store flag not registered on envset command")
	}
}
