package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envoy-cli/internal/alias"
)

var aliasStorePath = filepath.Join(os.Getenv("HOME"), ".envoy", "aliases.json")

func init() {
	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage profile aliases",
		Long:  "Create, list, and remove short aliases that map to profile names.",
	}

	setCmd := &cobra.Command{
		Use:   "set <alias> <profile>",
		Short: "Create or update an alias",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := alias.NewStore(aliasStorePath)
			if err != nil {
				return err
			}
			if err := s.Set(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("alias %q -> %q saved\n", args[0], args[1])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := alias.NewStore(aliasStorePath)
			if err != nil {
				return err
			}
			all := s.List()
			if len(all) == 0 {
				fmt.Println("no aliases defined")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ALIAS\tPROFILE")
			for k, v := range all {
				fmt.Fprintf(w, "%s\t%s\n", k, v)
			}
			return w.Flush()
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := alias.NewStore(aliasStorePath)
			if err != nil {
				return err
			}
			if err := s.Remove(args[0]); err != nil {
				return err
			}
			fmt.Printf("alias %q removed\n", args[0])
			return nil
		},
	}

	aliasCmd.AddCommand(setCmd, listCmd, removeCmd)
	rootCmd.AddCommand(aliasCmd)
}
