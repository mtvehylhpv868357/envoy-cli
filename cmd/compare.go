package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/compare"
	"github.com/user/envoy-cli/internal/profile"
)

func init() {
	var storePath string

	compareCmd := &cobra.Command{
		Use:   "compare <profileA> <profileB>",
		Short: "Compare two environment profiles side-by-side",
		Long: `Display a table showing which variables are unique to each profile,
which differ in value, and which are identical.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if storePath == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return err
				}
				storePath = home + "/.config/envoy/profiles"
			}

			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("loading profiles: %w", err)
			}

			result, err := compare.Profiles(store, args[0], args[1])
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "KEY\tSTATUS\t%s\t%s\n", result.ProfileA, result.ProfileB)
			fmt.Fprintf(w, "---\t------\t%s\t%s\n", "-------", "-------")

			for _, key := range result.AllKeys() {
				switch {
				case result.Same[key] != "":
					fmt.Fprintf(w, "%s\t=\t%s\t%s\n", key, result.Same[key], result.Same[key])
				case result.Differ[key] != [2]string{}:
					p := result.Differ[key]
					fmt.Fprintf(w, "%s\t~\t%s\t%s\n", key, p[0], p[1])
				case result.OnlyInA[key] != "":
					fmt.Fprintf(w, "%s\t+A\t%s\t(missing)\n", key, result.OnlyInA[key])
				case result.OnlyInB[key] != "":
					fmt.Fprintf(w, "%s\t+B\t(missing)\t%s\n", key, result.OnlyInB[key])
				}
			}
			return w.Flush()
		},
	}

	compareCmd.Flags().StringVar(&storePath, "store", "", "path to profile store directory")
	rootCmd.AddCommand(compareCmd)
}
