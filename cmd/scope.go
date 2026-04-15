package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-cli/internal/scope"
)

func scopeStorePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "envoy", "scope.json")
}

func init() {
	scopeCmd := &cobra.Command{
		Use:   "scope",
		Short: "Bind profiles to directories",
		Long:  "Manage automatic profile activation based on the current directory.",
	}

	// scope set <dir> <profile>
	setCmd := &cobra.Command{
		Use:   "set <directory> <profile>",
		Short: "Bind a directory to a profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := scope.NewStore(scopeStorePath())
			if err != nil {
				return err
			}
			dir, _ := filepath.Abs(args[0])
			if err := st.Set(dir, args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Bound %q → %q\n", dir, args[1])
			return nil
		},
	}

	// scope get [dir]
	getCmd := &cobra.Command{
		Use:   "get [directory]",
		Short: "Show the profile bound to a directory (walks up if needed)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}
			dir, _ = filepath.Abs(dir)
			st, err := scope.NewStore(scopeStorePath())
			if err != nil {
				return err
			}
			p := st.Resolve(dir)
			if p == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "No profile bound for this directory.")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), p)
			return nil
		},
	}

	// scope unset <dir>
	unsetCmd := &cobra.Command{
		Use:   "unset <directory>",
		Short: "Remove the binding for a directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := scope.NewStore(scopeStorePath())
			if err != nil {
				return err
			}
			dir, _ := filepath.Abs(args[0])
			return st.Remove(dir)
		},
	}

	// scope list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all directory-to-profile bindings",
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := scope.NewStore(scopeStorePath())
			if err != nil {
				return err
			}
			bindings := st.List()
			if len(bindings) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No scope bindings defined.")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "DIRECTORY\tPROFILE")
			for _, b := range bindings {
				fmt.Fprintf(w, "%s\t%s\n", b.Directory, b.Profile)
			}
			return w.Flush()
		},
	}

	scopeCmd.AddCommand(setCmd, getCmd, unsetCmd, listCmd)
	rootCmd.AddCommand(scopeCmd)
}
