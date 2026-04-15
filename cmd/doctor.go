package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/doctor"
	"github.com/yourorg/envoy-cli/internal/profile"
)

func init() {
	var storePath string

	doctorCmd := &cobra.Command{
		Use:   "doctor [profile]",
		Short: "Run health checks on an environment profile",
		Long: `Diagnose potential issues in an environment profile.

Checks include empty values, lowercase keys, unresolved variable references,
and secrets with suspiciously short values.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("load store: %w", err)
			}

			p, err := store.Get(name)
			if err != nil {
				return fmt.Errorf("profile %q not found: %w", name, err)
			}

			report := doctor.Check(p.Vars)

			if len(report.Findings) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "✓ No issues found in profile %q\n", name)
				return nil
			}

			for _, f := range report.Findings {
				var icon string
				switch f.Severity {
				case doctor.SeverityError:
					icon = "✗"
				case doctor.SeverityWarning:
					icon = "!"
				default:
					icon = "·"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s [%s] %s: %s\n", icon, f.Severity, f.Key, f.Message)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", report.Summary())

			if report.HasErrors() {
				os.Exit(1)
			}
			return nil
		},
	}

	doctorCmd.Flags().StringVar(&storePath, "store", profile.DefaultStorePath(), "path to profile store directory")

	rootCmd.AddCommand(doctorCmd)
}
