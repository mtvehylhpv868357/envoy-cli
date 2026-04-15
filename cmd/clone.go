package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/clone"
	"github.com/user/envoy-cli/internal/profile"
)

func init() {
	var (
		overrides []string
		setActive bool
	)

	cloneCmd := &cobra.Command{
		Use:   "clone <source> <destination>",
		Short: "Duplicate a profile under a new name",
		Long: `Clone copies an existing profile into a new profile.
You can override specific variables in the copy using --set KEY=VALUE flags.
The original profile is left unchanged.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]

			store, err := profile.LoadStore(profileStorePath())
			if err != nil {
				return fmt.Errorf("loading profile store: %w", err)
			}

			opts := clone.DefaultOptions()
			opts.SetActive = setActive

			for _, kv := range overrides {
				parts := strings.SplitN(kv, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid --set value %q: expected KEY=VALUE", kv)
				}
				opts.Overrides[parts[0]] = parts[1]
			}

			if err := clone.Profile(store, src, dst, opts); err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Profile %q cloned to %q\n", src, dst)
			if setActive {
				fmt.Fprintf(os.Stdout, "Active profile set to %q\n", dst)
			}
			return nil
		},
	}

	cloneCmd.Flags().StringArrayVar(&overrides, "set", nil, "Override a variable in the clone (KEY=VALUE, repeatable)")
	cloneCmd.Flags().BoolVar(&setActive, "activate", false, "Set the cloned profile as the active profile")

	rootCmd.AddCommand(cloneCmd)
}
