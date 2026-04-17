package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envoy-cli/internal/filter"
	"envoy-cli/internal/profile"
)

func init() {
	var keyPattern, valuePattern, prefix string
	var invert bool

	filterCmd := &cobra.Command{
		Use:   "filter <profile>",
		Short: "Filter environment variables in a profile by key, value, or prefix",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := profile.LoadStore("")
			if err != nil {
				return fmt.Errorf("load store: %w", err)
			}
			p, err := store.Get(args[0])
			if err != nil {
				return fmt.Errorf("profile %q not found: %w", args[0], err)
			}

			opts := filter.DefaultOptions()
			opts.KeyPattern = keyPattern
			opts.ValuePattern = valuePattern
			opts.Prefix = prefix
			opts.Invert = invert

			result, err := filter.Map(p.Vars, opts)
			if err != nil {
				return fmt.Errorf("filter: %w", err)
			}

			keys := make([]string, 0, len(result))
			for k := range result {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			w := cmd.OutOrStdout()
			if len(keys) == 0 {
				fmt.Fprintln(w, "no matching variables")
				return nil
			}
			for _, k := range keys {
				fmt.Fprintf(w, "%s=%s\n", k, result[k])
			}
			return nil
		},
	}

	filterCmd.Flags().StringVar(&keyPattern, "key", "", "regex pattern to match keys")
	filterCmd.Flags().StringVar(&valuePattern, "value", "", "regex pattern to match values")
	filterCmd.Flags().StringVar(&prefix, "prefix", "", "filter by key prefix")
	filterCmd.Flags().BoolVar(&invert, "invert", false, "invert the filter (exclude matches)")

	_ = os.Getenv // suppress unused import
	rootCmd.AddCommand(filterCmd)
}
