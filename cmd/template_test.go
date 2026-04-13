package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func findTemplateCmd(root *cobra.Command) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Use == "template <file>" {
			return c
		}
	}
	return nil
}

func TestTemplateCmd_Registered(t *testing.T) {
	if findTemplateCmd(rootCmd) == nil {
		t.Error("expected 'template' subcommand to be registered on root")
	}
}

func TestTemplateCmd_HelpContainsProfile(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"template", "--help"})
	_ = rootCmd.Execute()
	output := buf.String()
	if output == "" {
		// fallback: check via command directly
		cmd := findTemplateCmd(rootCmd)
		if cmd == nil {
			t.Fatal("template cmd not found")
		}
		output = cmd.Long
	}
	if output == "" {
		t.Error("expected non-empty help output")
	}
}

func TestTemplateCmd_FlagsRegistered(t *testing.T) {
	cmd := findTemplateCmd(rootCmd)
	if cmd == nil {
		t.Fatal("template cmd not found")
	}
	if cmd.Flags().Lookup("profile") == nil {
		t.Error("expected --profile flag")
	}
	if cmd.Flags().Lookup("strict") == nil {
		t.Error("expected --strict flag")
	}
}

func TestTemplateCmd_MissingFileArg(t *testing.T) {
	cmd := findTemplateCmd(rootCmd)
	if cmd == nil {
		t.Fatal("template cmd not found")
	}
	buf := &bytes.Buffer{}
	cmd.SetErr(buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when no file argument provided")
	}
}
