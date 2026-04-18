package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/envoy-cli/internal/profile"
	"github.com/envoy-cli/internal/split"
	"github.com/spf13/cobra"
)

func init() {
	var storePath string
	var overwrite bool
	var keepPrefix bool

	cmd := &cobra.Command{
		Use:   "split <profile> <PREFIX1> [PREFIX2 ...]",
		Short: "Split a profile into multiple profiles by key prefix",
		Long: `Split a single profile into several smaller profiles, one per prefix.

Each prefix (e.g. DB_, APP_) becomes its own profile named after the
lower-cased prefix. By default the prefix is stripped from the key names
in the resulting profiles.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			source := args[0]
			prefixes := args[1:]

			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("loading store: %w", err)
			}

			opts := split.DefaultOptions()
			opts.Overwrite = overwrite
			opts.StripPrefix = !keepPrefix

			created, err := split.ByPrefix(store, source, prefixes, opts)
			if err != nil {
				return err
			}

			if len(created) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No profiles created (no keys matched the given prefixes).")
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created profiles: %s\n", strings.Join(created, ", "))
			return nil
		},
	}

	defaultStore := filepath.Join(defaultProfileDir(), "profiles")
	cmd.Flags().StringVar(&storePath, "store", defaultStore, "path to profile store directory")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing destination profiles")
	cmd.Flags().BoolVar(&keepPrefix, "keep-prefix", false, "do not strip the prefix from destination keys")

	rootCmd.AddCommand(cmd)
}
