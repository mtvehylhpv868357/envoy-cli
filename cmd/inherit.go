package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/envoy/internal/inherit"
	"github.com/envoy-cli/envoy/internal/profile"
	"github.com/spf13/cobra"
)

func init() {
	var (
		storePath  string
		overwrite  bool
		strict     bool
		setFlags   []string
	)

	cmd := &cobra.Command{
		Use:   "inherit <base> <child>",
		Short: "Create a profile that inherits variables from a base profile",
		Long: `Create a new profile <child> whose variables are inherited from <base>.

Optional --set KEY=VALUE flags override or extend inherited variables.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			base, dst := args[0], args[1]

			st, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("loading profile store: %w", err)
			}

			overrides := make(map[string]string, len(setFlags))
			for _, pair := range setFlags {
				parts := strings.SplitN(pair, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid --set value %q: expected KEY=VALUE", pair)
				}
				overrides[parts[0]] = parts[1]
			}

			opts := inherit.DefaultOptions()
			opts.Overwrite = overwrite
			opts.Strict = strict

			if err := inherit.Profile(st, base, dst, overrides, opts); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Profile %q created from base %q\n", dst, base)
			return nil
		},
	}

	defaultStore := os.Getenv("ENVOY_PROFILE_DIR")
	if defaultStore == "" {
		defaultStore = defaultProfileDir()
	}

	cmd.Flags().StringVar(&storePath, "store", defaultStore, "path to profile store directory")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination profile if it already exists")
	cmd.Flags().BoolVar(&strict, "strict", true, "error if base profile does not exist")
	cmd.Flags().StringArrayVar(&setFlags, "set", nil, "override or add variables (KEY=VALUE, repeatable)")

	rootCmd.AddCommand(cmd)
}
