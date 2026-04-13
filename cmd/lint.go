package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
	"envoy-cli/internal/lint"
	"envoy-cli/internal/profile"
)

func init() {
	var fileFlag string
	var profileFlag string

	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Check environment variables for common issues",
		Long:  "Lint validates environment variable keys and values against common rules such as naming conventions and secret exposure.",
		RunE: func(cmd *cobra.Command, args []string) error {
			vars := map[string]string{}

			if fileFlag != "" {
				loaded, err := env.LoadFromFile(fileFlag)
				if err != nil {
					return fmt.Errorf("failed to load file %q: %w", fileFlag, err)
				}
				vars = loaded
			} else if profileFlag != "" {
				store, err := profile.LoadStore("")
				if err != nil {
					return fmt.Errorf("failed to load profile store: %w", err)
				}
				p, err := store.Get(profileFlag)
				if err != nil {
					return fmt.Errorf("profile %q not found: %w", profileFlag, err)
				}
				vars = p.Vars
			} else {
				return fmt.Errorf("provide --file or --profile to lint")
			}

			report := lint.Check(vars)
			for _, issue := range report.Issues {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s: %s\n", issue.Severity, issue.Key, issue.Message)
			}
			fmt.Fprintln(cmd.OutOrStdout(), report.Summary())

			if report.HasErrors() {
				os.Exit(1)
			}
			return nil
		},
	}

	lintCmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Path to a .env file to lint")
	lintCmd.Flags().StringVarP(&profileFlag, "profile", "p", "", "Profile name to lint")

	rootCmd.AddCommand(lintCmd)
}
