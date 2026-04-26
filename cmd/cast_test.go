package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func findCastCmd(t *testing.T) bool {
	t.Helper()
	for _, c := range rootCmd.Commands() {
		if c.Use == "cast <profile> <type>" {
			return true
		}
	}
	return false
}

func TestCastCmd_Registered(t *testing.T) {
	if !findCastCmd(t) {
		t.Fatal("cast command not registered on root")
	}
}

func TestCastCmd_RequiresTwoArgs(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"cast", "myprofile"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when second arg (type) is missing")
	}
}

func TestCastCmd_HelpContainsTypes(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"cast", "--help"})
	_ = rootCmd.Execute()
	if !strings.Contains(buf.String(), "bool") || !strings.Contains(buf.String(), "duration") {
		t.Errorf("help text should mention supported types, got:\n%s", buf.String())
	}
}

func TestCastCmd_StoreFlagRegistered(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "cast <profile> <type>" {
			if f := c.Flags().Lookup("store"); f == nil {
				t.Fatal("expected --store flag on cast command")
			}
			return
		}
	}
	t.Fatal("cast command not found")
}

func TestCastValue_UnknownType(t *testing.T) {
	_, err := castValue("hello", "xml")
	if err == nil || !strings.Contains(err.Error(), "unknown type") {
		t.Fatalf("expected unknown type error, got %v", err)
	}
}

func TestCastValue_Bool(t *testing.T) {
	v, err := castValue("yes", "bool")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != true {
		t.Errorf("expected true, got %v", v)
	}
}

func TestCastValue_Int(t *testing.T) {
	v, err := castValue("99", "int")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 99 {
		t.Errorf("expected 99, got %v", v)
	}
}
