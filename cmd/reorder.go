package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/envoy-cli/internal/profile"
	"github.com/envoy-cli/internal/reorder"
	"github.com/spf13/cobra"
)

func init() {
	var (
		storeDir  string
		strategy string
		customOrder string
	)

	cmd := &cobra.Command{
		Use:   "reorder <profile>",
		Short: "Reorder keys in a profile",
		Long: `Reorder the keys of an environment profile.

Strategies:
  alpha       ascending alphabetical (default)
  alpha-desc  descending alphabetical
  custom      order defined by --order flag; unlisted keys follow alphabetically`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := profile.LoadStore(storeDir)
			if err != nil {
				return fmt.Errorf("loading store: %w", err)
			}

			opts := reorder.Options{
				Strategy: reorder.Strategy(strategy),
			}
			if customOrder != "" {
				opts.Order = strings.Split(customOrder, ",")
			}

			vars, err := reorder.Profile(st, args[0], opts)
			if err != nil {
				return err
			}

			keys := reorder.Keys(vars, opts)
			for _, k := range keys {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, vars[k])
			}
			return nil
		},
	}

	defaultStore := filepath.Join(os.Getenv("HOME"), ".config", "envoy", "profiles")
	cmd.Flags().StringVar(&storeDir, "store", defaultStore, "profile store directory")
	cmd.Flags().StringVar(&strategy, "strategy", string(reorder.StrategyAlpha), "reorder strategy (alpha|alpha-desc|custom)")
	cmd.Flags().StringVar(&customOrder, "order", "", "comma-separated key order for custom strategy")

	rootCmd.AddCommand(cmd)
}
