package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func findPatchCmd(t *testing.T) bool {
	t.Helper()
	for _, c := range rootCmd.Commands() {
		if c.Name() == "patch" {
			return true
		}
	}
	return false
}

func TestPatchCmd_Registered(t *testing.T) {
	if !findPatchCmd(t) {
		t.Fatal("patch command not registered")
	}
}

func TestPatchCmd_RequiresArg(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"patch"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no profile name given")
	}
}

func TestPatchCmd_HelpContainsPatch(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"patch", "--help"})
	_ = rootCmd.Execute()
	if !strings.Contains(buf.String(), "patch") {
		t.Error("help output should mention patch")
	}
}

func TestPatchCmd_DeleteFlagRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "patch" {
			if f := c.Flags().Lookup("delete"); f == nil {
				t.Error("expected --delete flag")
			}
			return
		}
	}
	t.Fatal("patch command not found")
}

func TestPatchCmd_StoreFlagRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "patch" {
			if f := c.Flags().Lookup("store"); f == nil {
				t.Error("expected --store flag")
			}
			return
		}
	}
	t.Fatal("patch command not found")
}
