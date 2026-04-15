package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/pin"
)

var pinStorePath = filepath.Join(os.Getenv("HOME"), ".config", "envoy", "pins.json")

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Pin a profile to a directory for automatic activation",
	}

	setCmd := &cobra.Command{
		Use:   "set <profile>",
		Short: "Pin the current directory to a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			s, err := pin.NewStore(pinStorePath)
			if err != nil {
				return err
			}
			if err := s.Set(cwd, args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Pinned %q → %s\n", cwd, args[0])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Show the profile pinned to the current directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			s, err := pin.NewStore(pinStorePath)
			if err != nil {
				return err
			}
			name, err := s.Get(cwd)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), name)
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove the pin for the current directory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			s, err := pin.NewStore(pinStorePath)
			if err != nil {
				return err
			}
			return s.Remove(cwd)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all directory pins",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := pin.NewStore(pinStorePath)
			if err != nil {
				return err
			}
			pins := s.List()
			if len(pins) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No pins configured.")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "DIRECTORY\tPROFILE")
			for dir, profile := range pins {
				fmt.Fprintf(w, "%s\t%s\n", dir, profile)
			}
			return w.Flush()
		},
	}

	pinCmd.AddCommand(setCmd, getCmd, removeCmd, listCmd)
	rootCmd.AddCommand(pinCmd)
}
