package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/envoy-cli/internal/patch"
	"github.com/envoy-cli/internal/profile"
	"github.com/spf13/cobra"
)

func init() {
	var deleteKeys []string
	var storePath string

	cmd := &cobra.Command{
		Use:   "patch <profile> [KEY=VALUE ...]",
		Short: "Apply partial updates to a profile",
		Long:  "Add, overwrite, or delete individual keys in an existing profile without replacing it entirely.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if storePath == "" {
				home, _ := os.UserHomeDir()
				storePath = home + "/.config/envoy/profiles"
			}

			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("load store: %w", err)
			}

			opts := patch.DefaultOptions()
			opts.Delete = deleteKeys

			for _, pair := range args[1:] {
				parts := strings.SplitN(pair, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid key=value pair: %q", pair)
				}
				opts.Upsert[parts[0]] = parts[1]
			}

			result, err := patch.Profile(store, name, opts)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "patched profile %q (%d keys)\n", name, len(result))
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&deleteKeys, "delete", "d", nil, "keys to remove from the profile")
	cmd.Flags().StringVar(&storePath, "store", "", "path to profile store directory")

	rootCmd.AddCommand(cmd)
}
