package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/archive"
	"github.com/yourorg/envoy-cli/internal/profile"
)

func init() {
	var storeDir string
	var overwrite bool

	archiveCmd := &cobra.Command{
		Use:   "archive",
		Short: "Pack and unpack profile archives",
	}

	packCmd := &cobra.Command{
		Use:   "pack [output.tar.gz] [profile...]",
		Short: "Pack profiles into a gzipped tar archive",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			outPath := args[0]
			profiles := args[1:]
			if storeDir == "" {
				home, _ := os.UserHomeDir()
				storeDir = filepath.Join(home, ".envoy", "profiles")
			}
			f, err := os.Create(outPath)
			if err != nil {
				return fmt.Errorf("cannot create output file: %w", err)
			}
			defer f.Close()
			if err := archive.Pack(storeDir, profiles, f); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "packed %d profile(s) to %s\n", len(profiles), outPath)
			return nil
		},
	}

	unpackCmd := &cobra.Command{
		Use:   "unpack [archive.tar.gz]",
		Short: "Unpack profiles from a gzipped tar archive into the store",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if storeDir == "" {
				home, _ := os.UserHomeDir()
				storeDir = filepath.Join(home, ".envoy", "profiles")
			}
			_ = profile.LoadStore(storeDir) // ensure dir exists
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("cannot open archive: %w", err)
			}
			defer f.Close()
			names, err := archive.Unpack(f, storeDir, overwrite)
			if err != nil {
				return err
			}
			for _, n := range names {
				fmt.Fprintf(cmd.OutOrStdout(), "restored: %s\n", n)
			}
			return nil
		},
	}

	packCmd.Flags().StringVar(&storeDir, "store", "", "profile store directory")
	unpackCmd.Flags().StringVar(&storeDir, "store", "", "profile store directory")
	unpackCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing profiles")

	archiveCmd.AddCommand(packCmd, unpackCmd)
	rootCmd.AddCommand(archiveCmd)
}
