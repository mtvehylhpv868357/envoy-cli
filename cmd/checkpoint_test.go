package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findCheckpointCmd(root *cobra.Command) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Name() == "checkpoint" {
			return c
		}
	}
	return nil
}

func TestCheckpointCmd_Registered(t *testing.T) {
	if findCheckpointCmd(rootCmd) == nil {
		t.Fatal("expected 'checkpoint' command to be registered")
	}
}

func TestCheckpointCmd_SubcommandsRegistered(t *testing.T) {
	cp := findCheckpointCmd(rootCmd)
	if cp == nil {
		t.Fatal("checkpoint command not found")
	}
	want := []string{"save", "restore", "list", "delete"}
	names := map[string]bool{}
	for _, sub := range cp.Commands() {
		names[sub.Name()] = true
	}
	for _, w := range want {
		if !names[w] {
			t.Errorf("expected subcommand %q to be registered", w)
		}
	}
}

func TestCheckpointCmd_SaveRequiresTwoArgs(t *testing.T) {
	cp := findCheckpointCmd(rootCmd)
	var save *cobra.Command
	for _, sub := range cp.Commands() {
		if sub.Name() == "save" {
			save = sub
			break
		}
	}
	if save == nil {
		t.Fatal("save subcommand not found")
	}
	buf := &bytes.Buffer{}
	save.SetErr(buf)
	save.SetOut(buf)
	save.SetArgs([]string{"only-one"})
	if err := save.Execute(); err == nil {
		t.Error("expected error when only one arg provided to save")
	}
}

func TestCheckpointCmd_HelpContainsProfile(t *testing.T) {
	cp := findCheckpointCmd(rootCmd)
	buf := &bytes.Buffer{}
	cp.SetOut(buf)
	cp.SetArgs([]string{"--help"})
	_ = cp.Execute()
	if !strings.Contains(buf.String(), "profile") {
		t.Error("expected help text to mention 'profile'")
	}
}

func TestCheckpointCmd_NoteFlag(t *testing.T) {
	cp := findCheckpointCmd(rootCmd)
	var save *cobra.Command
	for _, sub := range cp.Commands() {
		if sub.Name() == "save" {
			save = sub
			break
		}
	}
	if save == nil {
		t.Fatal("save subcommand not found")
	}
	if f := save.Flags().Lookup("note"); f == nil {
		t.Error("expected --note flag on save subcommand")
	}
}
