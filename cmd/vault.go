package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/encrypt"
	"envoy-cli/internal/env"
)

var vaultPassphrase string

func init() {
	vaultCmd := &cobra.Command{
		Use:   "vault",
		Short: "Encrypt and decrypt .env profile files",
	}

	encryptCmd := &cobra.Command{
		Use:   "encrypt <file>",
		Short: "Encrypt a .env file into a .vault file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]
			vars, err := env.LoadFromFile(src)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			plaintext := []byte(env.Export(vars))
			ciphertext, err := encrypt.Encrypt(plaintext, vaultPassphrase)
			if err != nil {
				return fmt.Errorf("encryption failed: %w", err)
			}

			dest := src + ".vault"
			if err := os.WriteFile(dest, ciphertext, 0600); err != nil {
				return fmt.Errorf("failed to write vault file: %w", err)
			}

			fmt.Printf("Encrypted to %s\n", filepath.Base(dest))
			return nil
		},
	}

	decryptCmd := &cobra.Command{
		Use:   "decrypt <file>",
		Short: "Decrypt a .vault file and print its contents",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]
			if !strings.HasSuffix(src, ".vault") {
				return fmt.Errorf("expected a .vault file, got: %s", filepath.Base(src))
			}

			ciphertext, err := os.ReadFile(src)
			if err != nil {
				return fmt.Errorf("failed to read vault file: %w", err)
			}

			plaintext, err := encrypt.Decrypt(ciphertext, vaultPassphrase)
			if err != nil {
				return err
			}

			fmt.Println(string(plaintext))
			return nil
		},
	}

	encryptCmd.Flags().StringVarP(&vaultPassphrase, "passphrase", "p", "", "Passphrase for encryption (required)")
	_ = encryptCmd.MarkFlagRequired("passphrase")
	decryptCmd.Flags().StringVarP(&vaultPassphrase, "passphrase", "p", "", "Passphrase for decryption (required)")
	_ = decryptCmd.MarkFlagRequired("passphrase")

	vaultCmd.AddCommand(encryptCmd, decryptCmd)
	rootCmd.AddCommand(vaultCmd)
}
