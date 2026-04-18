package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envoy-cli/internal/archive"
)

func findArchiveCmd(t *testing.T) *cobra.Command {
	t.Helper()
	for _, c := range rootCmd.Commands() {
		if c.Use == "archive" {
			return c
		}
	}
	t.Fatal("archive command not registered")
	return nil
}

func TestArchiveCmd_Registered(t *testing.T) {
	findArchiveCmd(t)
}

func TestArchiveCmd_SubcommandsRegistered(t *testing.T) {
	cmd := findArchiveCmd(t)
	names := map[string]bool{}
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}
	for _, want := range []string{"pack", "unpack"} {
		if !names[want] {
			t.Errorf("subcommand %q not registered", want)
		}
	}
}

func TestArchiveCmd_PackRequiresTwoArgs(t *testing.T) {
	root := rootCmd
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"archive", "pack", "only-one-arg"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error with fewer than 2 args")
	}
}

func TestArchiveCmd_UnpackRequiresArg(t *testing.T) {
	root := rootCmd
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs([]string{"archive", "unpack"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestArchiveCmd_PackUnpack_Integration(t *testing.T) {
	dir := t.TempDir()
	profileDir := filepath.Join(dir, "profiles")
	os.MkdirAll(profileDir, 0755)
	data, _ := json.Marshal(map[string]string{"KEY": "val"})
	os.WriteFile(filepath.Join(profileDir, "myprofile.json"), data, 0600)

	outFile := filepath.Join(dir, "out.tar.gz")
	var buf bytes.Buffer
	if err := archive.Pack(profileDir, []string{"myprofile"}, &buf); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(outFile, buf.Bytes(), 0600)

	dstDir := filepath.Join(dir, "dst")
	f, _ := os.Open(outFile)
	defer f.Close()
	names, err := archive.Unpack(f, dstDir, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 1 || !strings.Contains(names[0], "myprofile") {
		t.Errorf("unexpected names: %v", names)
	}
}
