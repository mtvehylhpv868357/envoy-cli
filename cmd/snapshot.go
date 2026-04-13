package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/profile"
	"github.com/yourorg/envoy-cli/internal/snapshot"
)

func snapshotStore() (*snapshot.Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return snapshot.NewStore(filepath.Join(home, ".envoy", "snapshots"))
}

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage environment profile snapshots",
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save the active profile as a named snapshot",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ps, err := profile.LoadStore()
		if err != nil {
			return err
		}
		active := ps.Active()
		if active == "" {
			return fmt.Errorf("no active profile; use 'envoy profile use <name>' first")
		}
		vars, err := ps.Get(active)
		if err != nil {
			return err
		}
		ss, err := snapshotStore()
		if err != nil {
			return err
		}
		snap := snapshot.Snapshot{Name: args[0], Profile: active, Vars: vars}
		if err := ss.Save(snap); err != nil {
			return err
		}
		fmt.Printf("Snapshot %q saved from profile %q\n", args[0], active)
		return nil
	},
}

var snapshotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved snapshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		ss, err := snapshotStore()
		if err != nil {
			return err
		}
		names, err := ss.List()
		if err != nil {
			return err
		}
		if len(names) == 0 {
			fmt.Println("No snapshots found.")
			return nil
		}
		fmt.Println(strings.Join(names, "\n"))
		return nil
	},
}

var snapshotDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a snapshot by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ss, err := snapshotStore()
		if err != nil {
			return err
		}
		if err := ss.Delete(args[0]); err != nil {
			return err
		}
		fmt.Printf("Snapshot %q deleted\n", args[0])
		return nil
	},
}

func init() {
	snapshotCmd.AddCommand(snapshotSaveCmd)
	snapshotCmd.AddCommand(snapshotListCmd)
	snapshotCmd.AddCommand(snapshotDeleteCmd)
	rootCmd.AddCommand(snapshotCmd)
}
