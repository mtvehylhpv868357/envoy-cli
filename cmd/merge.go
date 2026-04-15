package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/internal/merge"
	"github.com/envoy-cli/internal/profile"
	"github.com/spf13/cobra"
)

func init() {
	var storePath string
	var strategy string
	var overwrite bool

	mergeCmd := &cobra.Command{
		Use:   "merge <base-profile> <source-profile>",
		Short: "Merge variables from source profile into base profile",
		Long: `Merge environment variables from the source profile into the base profile.

Conflicting keys are resolved according to --strategy:
  ours   - keep the base profile value
  theirs - keep the source profile value (default)

Use --overwrite to persist the result back to the base profile.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			base, src := args[0], args[1]

			if storePath == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				storePath = home + "/.config/envoy/profiles"
			}

			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("failed to load profile store: %w", err)
			}

			var strat merge.Strategy
			switch strings.ToLower(strategy) {
			case "ours":
				strat = merge.StrategyOurs
			case "theirs", "":
				strat = merge.StrategyTheirs
			default:
				return fmt.Errorf("unknown strategy %q: use 'ours' or 'theirs'", strategy)
			}

			opts := merge.Options{Strategy: strat, Overwrite: overwrite}
			res, err := merge.Profiles(store, base, src, opts)
			if err != nil {
				return err
			}

			if len(res.Conflicts) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Conflicts resolved (%s): %s\n",
					strat, strings.Join(res.Conflicts, ", "))
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Merged %d variable(s) from %q into %q\n",
				len(res.Merged), src, base)

			if overwrite {
				fmt.Fprintf(cmd.OutOrStdout(), "Profile %q updated.\n", base)
			}
			return nil
		},
	}

	mergeCmd.Flags().StringVar(&storePath, "store", "", "path to profile store directory")
	mergeCmd.Flags().StringVar(&strategy, "strategy", "theirs", "conflict resolution strategy: ours|theirs")
	mergeCmd.Flags().BoolVar(&overwrite, "overwrite", false, "persist merged result back to base profile")

	rootCmd.AddCommand(mergeCmd)
}
