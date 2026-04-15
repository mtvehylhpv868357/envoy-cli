package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-cli/internal/mask"
	"github.com/yourusername/envoy-cli/internal/profile"
)

func init() {
	var storePath string
	var revealChars int
	var showAll bool

	maskCmd := &cobra.Command{
		Use:   "mask <profile>",
		Short: "Display a profile's variables with sensitive values masked",
		Long: `Print all environment variables in a profile.
Sensitive keys (containing SECRET, TOKEN, PASSWORD, KEY, etc.) are masked
by default, showing only the last few characters of the value.`,
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

			opts := mask.DefaultOptions()
			opts.RevealChars = revealChars

			vars := p.Vars
			if !showAll {
				vars = opts.Vars(p.Vars)
			}

			keys := make([]string, 0, len(vars))
			for k := range vars {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, vars[k])
			}
			return nil
		},
	}

	maskCmd.Flags().StringVar(&storePath, "store", defaultProfileDir(), "path to profile store directory")
	maskCmd.Flags().IntVar(&revealChars, "reveal", 4, "number of trailing characters to reveal for masked values")
	maskCmd.Flags().BoolVar(&showAll, "show-all", false, "disable masking and show all values in plain text")

	rootCmd.AddCommand(maskCmd)
}
