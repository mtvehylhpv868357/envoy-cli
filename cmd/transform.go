package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/profile"
	"envoy-cli/internal/transform"
)

func init() {
	var op string
	var keys []string
	var store string

	transformCmd := &cobra.Command{
		Use:   "transform <profile>",
		Short: "Apply a transformation to env var values in a profile",
		Long: `Transform applies an operation (uppercase, lowercase, trim, base64encode, base64decode)
to all or selected keys in a profile and saves the result.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			s, err := profile.LoadStore(store)
			if err != nil {
				return fmt.Errorf("load store: %w", err)
			}
			p, err := s.Get(name)
			if err != nil {
				return fmt.Errorf("get profile %q: %w", name, err)
			}
			opts := transform.Options{
				Op:   transform.Op(op),
				Keys: keys,
			}
			result, err := transform.Map(p.Vars, opts)
			if err != nil {
				return fmt.Errorf("transform: %w", err)
			}
			p.Vars = result
			if err := s.Add(name, p.Vars); err != nil {
				return fmt.Errorf("save profile: %w", err)
			}
			fmt.Fprintf(os.Stdout, "transformed profile %q with op=%s\n", name, op)
			return nil
		},
	}

	transformCmd.Flags().StringVarP(&op, "op", "o", "trim", "Operation to apply (uppercase|lowercase|trim|base64encode|base64decode)")
	transformCmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "Comma-separated keys to target (default: all)")
	transformCmd.Flags().StringVarP(&store, "store", "s", "", "Path to profile store directory")

	rootCmd.AddCommand(transformCmd)
}
