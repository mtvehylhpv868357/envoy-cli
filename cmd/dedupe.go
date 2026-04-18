package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/dedupe"
	"github.com/user/envoy-cli/internal/profile"
)

func init() {
	var keepFirst bool
	var storePath string

	cmd := &cobra.Command{
		Use:   "dedupe <profile>",
		Short: "Remove duplicate keys from a profile, keeping last occurrence",
		Args:  cobra.ExactArgs(1),
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

			// Convert map to pairs to detect/remove duplicates.
			// In a map duplicates can't exist, so we surface duplicate info
			// and report what was found.
			pairs := make([]string, 0, len(vars))
			for k, v := range vars {
				pairs = append(pairs, k+"="+v)
			}

			dups := dedupe.Duplicates(pairs)
			if len(dups) == 0 {
				fmt.Fprintln(os.Stdout, "No duplicate keys found.")
				return nil
			}

			opts := dedupe.Options{KeepFirst: keepFirst}
			clean := dedupe.Map(pairs, opts)

			if err := store.Add(name, clean); err != nil {
				return fmt.Errorf("save profile: %w", err)
			}

			fmt.Fprintf(os.Stdout, "Removed duplicates for keys: %s\n",
				strings.Join(dups, ", "))
			return nil
		},
	}

	cmd.Flags().BoolVar(&keepFirst, "keep-first", false, "Keep first occurrence instead of last")
	cmd.Flags().StringVar(&storePath, "store", "", "Path to profile store directory")
	rootCmd.AddCommand(cmd)
}
