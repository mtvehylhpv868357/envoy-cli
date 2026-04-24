package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/profile"
	"envoy-cli/internal/stats"
)

func init() {
	var storePath string
	var topN int

	statsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Show aggregated statistics across environment profiles",
		Long: `Compute and display statistics for all profiles in the store.

Includes total variable counts, shared vs unique keys, empty values,
and the most frequently occurring keys across profiles.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if storePath == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("could not determine home directory: %w", err)
				}
				storePath = filepath.Join(home, ".envoy", "profiles")
			}

			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("failed to load profile store: %w", err)
			}

			names := store.List()
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No profiles found.")
				return nil
			}

			profiles := make(map[string]map[string]string, len(names))
			for _, name := range names {
				p, err := store.Get(name)
				if err != nil {
					continue
				}
				profiles[name] = p
			}

			r := stats.Compute(profiles)

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Profiles:       %d\n", r.ProfileCount)
			fmt.Fprintf(out, "Total vars:     %d\n", r.TotalVars)
			fmt.Fprintf(out, "Unique keys:    %d\n", r.UniqueKeys)
			fmt.Fprintf(out, "Shared keys:    %d\n", r.SharedKeys)
			fmt.Fprintf(out, "Empty values:   %d\n", r.EmptyValues)
			fmt.Fprintf(out, "Sensitive keys: %d\n", r.SensitiveKeys)

			if topN > 0 {
				top := stats.TopKeys(r, topN)
				if len(top) > 0 {
					fmt.Fprintf(out, "Top %d keys:    %s\n", topN, strings.Join(top, ", "))
				}
			}

			return nil
		},
	}

	statsCmd.Flags().StringVar(&storePath, "store", "", "Path to profile store directory")
	statsCmd.Flags().IntVar(&topN, "top", 5, "Number of top keys to display (0 to disable)")

	rootCmd.AddCommand(statsCmd)
}
