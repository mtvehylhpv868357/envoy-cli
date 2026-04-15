package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/tag"
)

func tagStorePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", "tags.json")
}

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags for environment profiles",
		Long:  "Add, remove, and filter profiles by tags to organise your environments.",
	}

	addCmd := &cobra.Command{
		Use:   "add <profile> <tag>",
		Short: "Add a tag to a profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := tag.NewStore(tagStorePath())
			if err != nil {
				return err
			}
			if err := s.Add(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Tagged profile %q with %q\n", args[0], args[1])
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <profile> <tag>",
		Short: "Remove a tag from a profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := tag.NewStore(tagStorePath())
			if err != nil {
				return err
			}
			if err := s.Remove(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed tag %q from profile %q\n", args[1], args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <tag>",
		Short: "List profiles with a given tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := tag.NewStore(tagStorePath())
			if err != nil {
				return err
			}
			profiles := s.FindByTag(args[0])
			if len(profiles) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No profiles tagged with %q\n", args[0])
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(profiles, "\n"))
			return nil
		},
	}

	tagCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(tagCmd)
}
