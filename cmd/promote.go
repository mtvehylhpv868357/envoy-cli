package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-cli/internal/promote"
)

func init() {
	var storePath string
	var overwrite bool
	var activate bool

	promoteCmd := &cobra.Command{
		Use:   "promote <source> <destination>",
		Short: "Promote a profile to another name (e.g. staging → production)",
		Long: `Promote copies all variables from the source profile into a new
destination profile. Optionally activate the destination after promotion.

Example:
  envoy promote staging production --activate`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]

			if storePath == "" {
				var err error
				storePath, err = defaultPromoteStorePath()
				if err != nil {
					return err
				}
			}

			opts := promote.DefaultOptions()
			opts.StorePath = storePath
			opts.Overwrite = overwrite
			opts.Activate = activate

			res, err := promote.Profile(src, dst, opts)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Promoted %q → %q\n", res.Source, res.Destination)
			if res.Activated {
				fmt.Fprintf(cmd.OutOrStdout(), "Active profile set to %q\n", res.Destination)
			}
			return nil
		},
	}

	promoteCmd.Flags().StringVar(&storePath, "store", "", "path to profile store directory")
	promoteCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination if it already exists")
	promoteCmd.Flags().BoolVarP(&activate, "activate", "a", false, "set destination as the active profile after promotion")

	rootCmd.AddCommand(promoteCmd)
}

func defaultPromoteStorePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("promote: resolve home dir: %w", err)
	}
	return home + "/.envoy/profiles", nil
}
