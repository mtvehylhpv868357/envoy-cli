package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func findInheritCmd(root *cobra.Command) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Name() == "inherit" {
			return c
		}
	}
	return nil
}

func TestInheritCmd_Registered(t *testing.T) {
	if findInheritCmd(rootCmd) == nil {
		t.Fatal("expected 'inherit' command to be registered")
	}
}

func TestInheritCmd_RequiresTwoArgs(t *testing.T) {
	cmd := findInheritCmd(rootCmd)
	if cmd == nil {
		t.Skip("inherit command not registered")
	}
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.RunE(cmd, []string{"only-one"})
	if err == nil {
		t.Fatal("expected error with fewer than two args")
	}
}

func TestInheritCmd_FlagsRegistered(t *testing.T) {
	cmd := findInheritCmd(rootCmd)
	if cmd == nil {
		t.Skip("inherit command not registered")
	}
	for _, flag := range []string{"store", "overwrite", "strict", "set"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Errorf("expected flag --%s to be registered", flag)
		}
	}
}

func TestInheritCmd_HelpContainsBase(t *testing.T) {
	cmd := findInheritCmd(rootCmd)
	if cmd == nil {
		t.Skip("inherit command not registered")
	}
	if cmd.Long == "" {
		t.Fatal("expected non-empty long description")
	}
	if !bytes.Contains([]byte(cmd.Long), []byte("base")) {
		t.Error("expected long description to mention 'base'")
	}
}

func TestInheritCmd_InvalidSetFlag(t *testing.T) {
	cmd := findInheritCmd(rootCmd)
	if cmd == nil {
		t.Skip("inherit command not registered")
	}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	// Manually invoke RunE with a bad --set value by setting the flag directly.
	if err := cmd.Flags().Set("set", "NOEQUALSIGN"); err != nil {
		t.Skipf("could not set flag: %v", err)
	}
	err := cmd.RunE(cmd, []string{"base", "child"})
	if err == nil {
		t.Error("expected error for malformed --set value")
	}
}
