package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-cli/internal/profile"
	"github.com/yourusername/envoy-cli/internal/rename"
)

func init() {
	var storePath string
	var overwrite bool

	renameCmd := &cobra.Command{
		Use:   "rename <old-name> <new-name>",
		Short: "Rename an environment profile",
		Long: `Rename an existing environment profile to a new name.

If the destination profile already exists the command will fail unless
--overwrite is supplied.  When the renamed profile is currently active the
active pointer is updated automatically.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName := args[0]
			newName := args[1]

			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("loading profile store: %w", err)
			}

			opts := rename.Options{Overwrite: overwrite}
			if err := rename.Profile(store, oldName, newName, opts); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}

			fmt.Printf("Profile %q renamed to %q\n", oldName, newName)
			return nil
		},
	}

	renameCmd.Flags().StringVar(&storePath, "store", defaultStorePath(), "path to profile store directory")
	renameCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination profile if it exists")

	rootCmd.AddCommand(renameCmd)
}
