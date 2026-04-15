package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/profile"
	"github.com/user/envoy-cli/internal/search"
)

func init() {
	var storeDir string
	var valuePattern string
	var exactKey bool

	cmd := &cobra.Command{
		Use:   "search [key-pattern]",
		Short: "Search for environment variables across all profiles",
		Long: `Search all profiles for keys or values matching the given pattern.

Examples:
  envoy search DB_HOST
  envoy search --value localhost
  envoy search DB_HOST --exact
  envoy search --value prod --key DB`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if storeDir == "" {
				storeDir = defaultProfileDir()
			}
			store, err := profile.LoadStore(storeDir)
			if err != nil {
				return fmt.Errorf("load store: %w", err)
			}

			opts := search.DefaultOptions()
			if len(args) > 0 {
				opts.KeyPattern = args[0]
			}
			opts.ValuePattern = valuePattern
			opts.ExactKey = exactKey

			results, err := search.Profiles(store, opts)
			if err != nil {
				return fmt.Errorf("search: %w", err)
			}

			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no matches found")
				return nil
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "PROFILE\tKEY\tVALUE")
			for _, r := range results {
				fmt.Fprintf(w, "%s\t%s\t%s\n", r.Profile, r.Key, r.Value)
			}
			return w.Flush()
		},
	}

	cmd.Flags().StringVar(&storeDir, "store", "", "profile store directory")
	cmd.Flags().StringVar(&valuePattern, "value", "", "filter by value substring")
	cmd.Flags().BoolVar(&exactKey, "exact", false, "require exact key match")

	rootCmd.AddCommand(cmd)
	_ = os.Getenv // suppress unused import
}

func defaultProfileDir() string {
	home, _ := os.UserHomeDir()
	return home + "/.envoy/profiles"
}
