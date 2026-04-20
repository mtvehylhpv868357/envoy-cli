package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/group"
)

var groupStorePath = filepath.Join(os.Getenv("HOME"), ".config", "envoy", "groups.json")

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage profile groups",
		Long:  "Organise profiles into named groups for bulk operations and discovery.",
	}

	// group add <group> <profile>
	addCmd := &cobra.Command{
		Use:   "add <group> <profile>",
		Short: "Add a profile to a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := group.NewStore(groupStorePath)
			if err != nil {
				return err
			}
			if err := s.Add(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added %q to group %q\n", args[1], args[0])
			return nil
		},
	}

	// group remove <group> <profile>
	removeCmd := &cobra.Command{
		Use:   "remove <group> <profile>",
		Short: "Remove a profile from a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := group.NewStore(groupStorePath)
			if err != nil {
				return err
			}
			return s.Remove(args[0], args[1])
		},
	}

	// group list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := group.NewStore(groupStorePath)
			if err != nil {
				return err
			}
			names := s.List()
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No groups defined.")
				return nil
			}
			for _, n := range names {
				members := s.Get(n)
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", n, strings.Join(members, ", "))
			}
			return nil
		},
	}

	// group delete <group>
	deleteCmd := &cobra.Command{
		Use:   "delete <group>",
		Short: "Delete an entire group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := group.NewStore(groupStorePath)
			if err != nil {
				return err
			}
			if err := s.Delete(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted group %q\n", args[0])
			return nil
		},
	}

	groupCmd.AddCommand(addCmd, removeCmd, listCmd, deleteCmd)
	rootCmd.AddCommand(groupCmd)
}
