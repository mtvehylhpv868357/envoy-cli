package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findCompareCmd(root *cobra.Command) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Name() == "compare" {
			return c
		}
	}
	return nil
}

func TestCompareCmd_Registered(t *testing.T) {
	if findCompareCmd(rootCmd) == nil {
		t.Error("compare command not registered")
	}
}

func TestCompareCmd_RequiresTwoArgs(t *testing.T) {
	cmd := findCompareCmd(rootCmd)
	if cmd == nil {
		t.Fatal("compare command not found")
	}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"only-one"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error with only one argument")
	}
}

func TestCompareCmd_HelpContainsProfiles(t *testing.T) {
	cmd := findCompareCmd(rootCmd)
	if cmd == nil {
		t.Fatal("compare command not found")
	}
	help := cmd.Long
	if !strings.Contains(help, "profile") {
		t.Errorf("expected help text to mention 'profile', got: %s", help)
	}
}

func TestCompareCmd_StoreFlagRegistered(t *testing.T) {
	cmd := findCompareCmd(rootCmd)
	if cmd == nil {
		t.Fatal("compare command not found")
	}
	if cmd.Flags().Lookup("store") == nil {
		t.Error("expected --store flag to be registered")
	}
}

func TestCompareCmd_ShortDescription(t *testing.T) {
	cmd := findCompareCmd(rootCmd)
	if cmd == nil {
		t.Fatal("compare command not found")
	}
	if !strings.Contains(cmd.Short, "Compare") {
		t.Errorf("expected Short to contain 'Compare', got: %s", cmd.Short)
	}
}
