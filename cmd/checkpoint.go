package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/your-org/envoy-cli/internal/checkpoint"
	"github.com/your-org/envoy-cli/internal/profile"
)

var checkpointStorePath = filepath.Join(os.Getenv("HOME"), ".envoy", "checkpoints")

func init() {
	store := func() *checkpoint.Store {
		s, err := checkpoint.NewStore(checkpointStorePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error opening checkpoint store:", err)
			os.Exit(1)
		}
		return s
	}

	cpCmd := &cobra.Command{
		Use:   "checkpoint",
		Short: "Manage named save-points for profile state",
	}

	saveCmd := &cobra.Command{
		Use:   "save <name> <profile>",
		Short: "Save the current state of a profile as a checkpoint",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, profileName := args[0], args[1]
			note, _ := cmd.Flags().GetString("note")
			dir, _ := cmd.Flags().GetString("store")
			ps, err := profile.LoadStore(dir)
			if err != nil {
				return err
			}
			vars, err := ps.Get(profileName)
			if err != nil {
				return err
			}
			return store().Save(checkpoint.Entry{
				Name: name, Profile: profileName, Vars: vars, Note: note,
			})
		},
	}
	saveCmd.Flags().String("note", "", "optional note to attach to the checkpoint")
	saveCmd.Flags().String("store", "", "profile store directory")

	restoreCmd := &cobra.Command{
		Use:   "restore <name> <target-profile>",
		Short: "Restore a checkpoint into a profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cpName, target := args[0], args[1]
			dir, _ := cmd.Flags().GetString("store")
			e, err := store().Load(cpName)
			if err != nil {
				return err
			}
			ps, err := profile.LoadStore(dir)
			if err != nil {
				return err
			}
			return ps.Add(target, e.Vars)
		},
	}
	restoreCmd.Flags().String("store", "", "profile store directory")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved checkpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := store().List()
			if err != nil {
				return err
			}
			if len(entries) == 0 {
				fmt.Println("no checkpoints found")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tPROFILE\tCREATED\tNOTE")
			for _, e := range entries {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.Name, e.Profile, e.CreatedAt.Format("2006-01-02 15:04"), e.Note)
			}
			return w.Flush()
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a checkpoint by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return store().Delete(args[0])
		},
	}

	cpCmd.AddCommand(saveCmd, restoreCmd, listCmd, deleteCmd)
	rootCmd.AddCommand(cpCmd)
}
