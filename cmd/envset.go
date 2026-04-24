package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/env"
	"github.com/user/envoy-cli/internal/envset"
	"github.com/user/envoy-cli/internal/profile"
)

func init() {
	var storePath string

	envsetCmd := &cobra.Command{
		Use:   "envset",
		Short: "Set-theoretic operations on environment profiles",
		Long:  "Perform union, intersection, difference, or symmetric-difference on two or more profiles.",
	}

	newSetOp := func(use, short string, op func(envset.Options, ...map[string]string) map[string]string) *cobra.Command {
		return &cobra.Command{
			Use:   use,
			Short: short,
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				store, err := profile.LoadStore(storePath)
				if err != nil {
					return err
				}
				opts := envset.DefaultOptions()
				maps := make([]map[string]string, 0, len(args))
				for _, name := range args {
					p, err := store.Get(name)
					if err != nil {
						return fmt.Errorf("profile %q not found: %w", name, err)
					}
					maps = append(maps, p.Vars)
				}
				result := op(opts, maps...)
				for _, line := range env.Export(result) {
					fmt.Fprintln(os.Stdout, line)
				}
				return nil
			},
		}
	}

	unionCmd := newSetOp("union <profile> [profile...]", "Output union of profiles", envset.Union)
	intersectCmd := newSetOp("intersect <profile> [profile...]", "Output intersection of profiles", envset.Intersect)

	diffCmd := &cobra.Command{
		Use:   "diff <profileA> <profileB>",
		Short: "Output keys in profileA not in profileB",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := profile.LoadStore(storePath)
			if err != nil {
				return err
			}
			opts := envset.DefaultOptions()
			a, err := store.Get(args[0])
			if err != nil {
				return fmt.Errorf("profile %q not found: %w", args[0], err)
			}
			b, err := store.Get(args[1])
			if err != nil {
				return fmt.Errorf("profile %q not found: %w", args[1], err)
			}
			result := envset.Difference(opts, a.Vars, b.Vars)
			for _, line := range env.Export(result) {
				fmt.Fprintln(os.Stdout, line)
			}
			return nil
		},
	}

	envsetCmd.PersistentFlags().StringVar(&storePath, "store", "", "Path to profile store directory")
	envsetCmd.AddCommand(unionCmd, intersectCmd, diffCmd)
	rootCmd.AddCommand(envsetCmd)
}
