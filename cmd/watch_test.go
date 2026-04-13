package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWatchCmd_Registered(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "watch [file]" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected 'watch' command to be registered")
	}
}

func TestWatchCmd_FlagRegistered(t *testing.T) {
	var watchCmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if strings.HasPrefix(c.Use, "watch") {
			watchCmd = c
			break
		}
	}
	if watchCmd == nil {
		t.Fatal("watch command not found")
	}
	if f := watchCmd.Flags().Lookup("interval"); f == nil {
		t.Fatal("expected --interval flag to be registered")
	}
}

func TestWatchCmd_MissingFile(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"watch", "/nonexistent/totally/missing.env"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	if !strings.Contains(err.Error(), "file not found") {
		t.Errorf("error = %q, want 'file not found'", err.Error())
	}
}

func TestWatchCmd_HelpContainsCtrlC(t *testing.T) {
	_ = os.MkdirTemp("", "") // ensure os import used
	_ = filepath.Abs(".")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"watch", "--help"})
	_ = rootCmd.Execute()
	if !strings.Contains(buf.String(), "Ctrl+C") {
		t.Errorf("help text missing 'Ctrl+C', got: %s", buf.String())
	}
}
