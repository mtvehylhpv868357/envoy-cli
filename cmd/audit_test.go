package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func findCmd(root *cobra.Command, name string) *cobra.Command {
	for _, sub := range root.Commands() {
		if sub.Name() == name {
			return sub
		}
	}
	return nil
}

func TestAuditCmd_SubcommandsRegistered(t *testing.T) {
	auditCmd := findCmd(rootCmd, "audit")
	if auditCmd == nil {
		t.Fatal("audit command not registered")
	}

	want := []string{"list", "clear"}
	for _, name := range want {
		found := false
		for _, sub := range auditCmd.Commands() {
			if sub.Name() == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("subcommand %q not registered under audit", name)
		}
	}
}

func TestAuditCmd_HelpText(t *testing.T) {
	auditCmd := findCmd(rootCmd, "audit")
	if auditCmd == nil {
		t.Fatal("audit command not registered")
	}
	buf := new(bytes.Buffer)
	auditCmd.SetOut(buf)
	auditCmd.SetArgs([]string{"--help"})
	_ = auditCmd.Execute()
	if !strings.Contains(buf.String(), "audit") {
		t.Errorf("help output missing 'audit': %q", buf.String())
	}
}

func TestAuditCmd_ListShortDescription(t *testing.T) {
	auditCmd := findCmd(rootCmd, "audit")
	if auditCmd == nil {
		t.Fatal("audit command not registered")
	}
	for _, sub := range auditCmd.Commands() {
		if sub.Name() == "list" {
			if sub.Short == "" {
				t.Error("audit list missing short description")
			}
			return
		}
	}
	t.Error("audit list subcommand not found")
}
