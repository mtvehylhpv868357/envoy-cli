package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-cli/internal/copy"
	"github.com/yourusername/envoy-cli/internal/profile"
)

func init() {
	var (
		storePath string
		overwrite bool
	)

	copyCmd := &cobra.Command{
		Use:   "copy <source> <destination>",
		Short: "Duplicate an existing profile under a new name",
		Long: `Copy duplicates a profile and saves it under a new name.

All environment variables from the source profile are copied verbatim.
Use --overwrite to replace an existing destination profile.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dst := args[0], args[1]

			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("failed to load profile store: %w", err)
			}

			opts := copy.Options{Overwrite: overwrite}
			if err := copy.Profile(store, src, dst, opts); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				return err
			}

			fmt.Printf("Profile %q copied to %q\n", src, dst)
			return nil
		},
	}

	copyCmd.Flags().StringVar(&storePath, "store", profile.DefaultStorePath(), "path to the profile store directory")
	copyCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite destination profile if it already exists")

	rootCmd.AddCommand(copyCmd)
}
