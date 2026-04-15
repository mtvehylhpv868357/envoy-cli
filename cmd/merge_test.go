package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func findMergeCmd(t *testing.T) bool {
	t.Helper()
	for _, c := range rootCmd.Commands() {
		if c.Name() == "merge" {
			return true
		}
	}
	return false
}

func TestMergeCmd_Registered(t *testing.T) {
	if !findMergeCmd(t) {
		t.Fatal("expected 'merge' command to be registered")
	}
}

func TestMergeCmd_RequiresTwoArgs(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"merge", "only-one"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when only one arg provided")
	}
}

func TestMergeCmd_HelpContainsProfiles(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"merge", "--help"})
	_ = rootCmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "profile") {
		t.Errorf("expected help text to mention 'profile', got:\n%s", out)
	}
}

func TestMergeCmd_FlagsRegistered(t *testing.T) {
	var mergeCmd = func() interface{ Flags() interface{ Lookup(string) interface{ Name string } } } {
		for _, c := range rootCmd.Commands() {
			if c.Name() == "merge" {
				return c
			}
		}
		return nil
	}()
	if mergeCmd == nil {
		t.Fatal("merge command not found")
	}

	for _, c := range rootCmd.Commands() {
		if c.Name() != "merge" {
			continue
		}
		for _, flag := range []string{"strategy", "overwrite", "store"} {
			if c.Flags().Lookup(flag) == nil {
				t.Errorf("expected flag --%s to be registered", flag)
			}
		}
	}
}

func TestMergeCmd_InvalidStrategy(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"merge", "--strategy", "invalid", "base", "src"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
	if !strings.Contains(err.Error(), "unknown strategy") {
		t.Errorf("expected 'unknown strategy' in error, got: %v", err)
	}
}
