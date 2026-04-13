package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVaultCmd_SubcommandsRegistered(t *testing.T) {
	vaultCmd, _, err := rootCmd.Find([]string{"vault"})
	if err != nil || vaultCmd == nil {
		t.Fatal("vault command not found")
	}

	names := map[string]bool{}
	for _, sub := range vaultCmd.Commands() {
		names[sub.Name()] = true
	}

	for _, expected := range []string{"encrypt", "decrypt"} {
		if !names[expected] {
			t.Errorf("expected subcommand %q to be registered", expected)
		}
	}
}

func TestVaultCmd_EncryptRequiresPassphrase(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"vault", "encrypt", "somefile.env"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when passphrase is missing")
	}
}

func TestVaultCmd_EncryptDecryptRoundTrip(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, "test.env")

	content := "DB_HOST=localhost\nDB_PORT=5432\n"
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	// Encrypt
	rootCmd.SetArgs([]string{"vault", "encrypt", "--passphrase", "secret123", envFile})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	vaultFile := envFile + ".vault"
	if _, err := os.Stat(vaultFile); os.IsNotExist(err) {
		t.Fatal("vault file was not created")
	}

	// Decrypt
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"vault", "decrypt", "--passphrase", "secret123", vaultFile})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "DB_HOST") {
		t.Errorf("expected decrypted output to contain DB_HOST, got: %s", output)
	}
}
