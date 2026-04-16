package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/normalize"
	"github.com/user/envoy-cli/internal/profile"
)

func init() {
	var (
		storePath   string
		dryRun     bool
		noUppercase bool
	)

	cmd := &cobra.Command{
		Use:   "normalize <profile>",
		Short: "Normalize environment variable keys in a profile",
		Long: `Normalize applies formatting rules to a profile's environment variable keys:
  - Uppercase all keys
  - Trim surrounding whitespace from keys and values
  - Replace hyphens with underscores in keys`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("load store: %w", err)
			}

			vars, err := store.Get(name)
			if err != nil {
				return fmt.Errorf("get profile %q: %w", name, err)
			}

			opts := normalize.DefaultOptions()
			if noUppercase {
				opts.UppercaseKeys = false
			}

			normalized, changed := normalize.Map(vars, opts)

			if len(changed) == 0 {
				fmt.Fprintln(os.Stdout, "profile already normalized, no changes needed")
				return nil
			}

			fmt.Fprintf(os.Stdout, "%d key(s) normalized:\n", len(changed))
			for _, k := range changed {
				fmt.Fprintf(os.Stdout, "  %s\n", k)
			}

			if dryRun {
				fmt.Fprintln(os.Stdout, "dry-run: no changes saved")
				return nil
			}

			if err := store.Add(name, normalized); err != nil {
				return fmt.Errorf("save profile: %w", err)
			}
			fmt.Fprintln(os.Stdout, "profile saved")
			return nil
		},
	}

	cmd.Flags().StringVar(&storePath, "store", defaultProfileDir(), "path to profile store")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without saving")
	cmd.Flags().BoolVar(&noUppercase, "no-uppercase", false, "skip uppercasing keys")

	rootCmd.AddCommand(cmd)
}
