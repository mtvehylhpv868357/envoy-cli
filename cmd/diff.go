package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/diff"
	"envoy-cli/internal/env"
	"envoy-cli/internal/profile"
)

func init() {
	var diffCmd = &cobra.Command{
		Use:   "diff <profile-a> <profile-b>",
		Short: "Show differences between two environment profiles",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := profile.LoadStore("")
			if err != nil {
				return fmt.Errorf("loading profiles: %w", err)
			}

			profA, err := store.Get(args[0])
			if err != nil {
				return fmt.Errorf("profile %q not found", args[0])
			}

			profB, err := store.Get(args[1])
			if err != nil {
				return fmt.Errorf("profile %q not found", args[1])
			}

			beforeVars := env.ParseDotEnv(profA.Vars)
			afterVars := env.ParseDotEnv(profB.Vars)

			changes := diff.Compare(beforeVars, afterVars)
			if len(changes) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No differences between %q and %q\n", args[0], args[1])
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Diff: %s → %s\n\n", args[0], args[1])
			for _, c := range changes {
				switch c.Action {
				case "added":
					fmt.Fprintf(cmd.OutOrStdout(), "  + %s=%s\n", c.Key, c.New)
				case "removed":
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s=%s\n", c.Key, c.Old)
				case "modified":
					fmt.Fprintf(cmd.OutOrStdout(), "  ~ %s: %s → %s\n", c.Key, c.Old, c.New)
				}
			}

			summary := diff.Summary(changes)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d added, %d removed, %d modified\n",
				summary["added"], summary["removed"], summary["modified"])
			return nil
		},
	}

	if rootCmd == nil {
		fmt.Fprintln(os.Stderr, "rootCmd is nil")
		os.Exit(1)
	}
	rootCmd.AddCommand(diffCmd)
}
