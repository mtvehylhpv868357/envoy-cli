package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/envsign"
	"github.com/yourorg/envoy-cli/internal/profile"
)

func init() {
	var passphrase string
	var storePath string
	var outputFile string

	signCmd := &cobra.Command{
		Use:   "sign <profile>",
		Short: "Sign a profile's vars with an HMAC signature",
		Long: `Sign creates a signed envelope for the given profile using HMAC-SHA256.
The envelope can later be verified with 'envoy verify' to detect tampering.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("load store: %w", err)
			}
			vars, err := store.Get(name)
			if err != nil {
				return fmt.Errorf("profile %q not found: %w", name, err)
			}
			envelope, err := envsign.Sign(name, vars, passphrase)
			if err != nil {
				return fmt.Errorf("sign: %w", err)
			}
			data, err := json.MarshalIndent(envelope, "", "  ")
			if err != nil {
				return err
			}
			if outputFile != "" {
				if err := os.WriteFile(outputFile, data, 0600); err != nil {
					return fmt.Errorf("write: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Signed envelope written to %s\n", outputFile)
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(data))
			return nil
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify <envelope-file>",
		Short: "Verify a signed envelope's integrity",
		Long:  `Verify reads a signed envelope JSON file and checks its HMAC signature.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("read: %w", err)
			}
			var envelope envsign.SignedEnvelope
			if err := json.Unmarshal(data, &envelope); err != nil {
				return fmt.Errorf("parse envelope: %w", err)
			}
			if err := envsign.Verify(&envelope, passphrase); err != nil {
				return fmt.Errorf("verification failed: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Signature valid ✓")
			return nil
		},
	}

	signCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase for HMAC signing (required)")
	signCmd.Flags().StringVar(&storePath, "store", "", "Path to profile store directory")
	signCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write envelope to file instead of stdout")
	_ = signCmd.MarkFlagRequired("passphrase")

	verifyCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase used when signing (required)")
	_ = verifyCmd.MarkFlagRequired("passphrase")

	rootCmd.AddCommand(signCmd)
	rootCmd.AddCommand(verifyCmd)
}
