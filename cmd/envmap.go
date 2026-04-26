package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/envmap"
	"github.com/yourorg/envoy-cli/internal/profile"
)

func init() {
	var storePath string
	var uppercase bool
	var noOverwrite bool

	cmd := &cobra.Command{
		Use:   "envmap <profile>",
		Short: "Convert a profile to KEY=VALUE environment pairs",
		Long: `Loads the named profile and prints its variables as a sorted list of
KEY=VALUE pairs, suitable for piping or further processing.

Use --uppercase to normalise all keys to uppercase before output.
Use --no-overwrite to see which keys would be skipped when merging with the
current OS environment.`,
		Args: cobra.ExactArgs(1),
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

			opts := envmap.Options{
				Overwrite:     !noOverwrite,
				UppercaseKeys: uppercase,
			}

			// Merge profile vars on top of (or alongside) the OS environment.
			base := envmap.FromOS()
			merged := envmap.Merge(base, vars, opts)

			// Only print the keys that came from the profile so the output
			// stays focused on what the profile contributes.
			profKeys := envmap.Keys(vars)
			for _, k := range profKeys {
				key := k
				if uppercase {
					key = strings.ToUpper(k)
				}
				v, ok := merged[key]
				if !ok {
					continue // skipped by no-overwrite
				}
				fmt.Fprintf(os.Stdout, "%s=%s\n", key, v)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&storePath, "store", profile.DefaultStorePath(), "path to profile store directory")
	cmd.Flags().BoolVar(&uppercase, "uppercase", false, "normalise all keys to uppercase")
	cmd.Flags().BoolVar(&noOverwrite, "no-overwrite", false, "skip keys already present in the OS environment")

	rootCmd.AddCommand(cmd)
}
