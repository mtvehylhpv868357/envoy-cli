package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/interpolate"
	"envoy-cli/internal/profile"
)

func init() {
	var storePath string
	var strict bool

	cmd := &cobra.Command{
		Use:   "interpolate <profile>",
		Short: "Interpolate variable references within a profile",
		Long: `Resolves $VAR, ${VAR}, and ${VAR:-default} references within a profile's
variables using other variables in the same profile or the current environment.`,
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
			opts := interpolate.DefaultOptions()
			opts.Strict = strict
			out, err := interpolate.Map(p.Vars, os.LookupEnv, opts)
			if err != nil {
				return err
			}
			for k, v := range out {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&storePath, "store", profile.DefaultStorePath(), "path to profile store")
	cmd.Flags().BoolVar(&strict, "strict", false, "error on unresolved variables")
	rootCmd.AddCommand(cmd)
}
