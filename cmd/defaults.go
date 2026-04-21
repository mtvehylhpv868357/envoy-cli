package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/defaults"
)

const defaultsStorePath = ".envoy/defaults"

func init() {
	defaultsCmd := &cobra.Command{
		Use:   "defaults",
		Short: "Manage default environment variable values",
		Long:  "Set, get, delete, and list default values applied to profiles when a key is absent.",
	}

	defaultsCmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a default value for a key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := defaults.NewStore(defaultsStorePath)
			if err != nil {
				return err
			}
			if err := s.Set(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "default set: %s=%s\n", args[0], args[1])
			return nil
		},
	})

	defaultsCmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get the default value for a key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := defaults.NewStore(defaultsStorePath)
			if err != nil {
				return err
			}
			v, ok := s.Get(args[0])
			if !ok {
				fmt.Fprintf(cmd.OutOrStdout(), "no default for %q])
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), v)
			return nil
		},
	})

	defaultsCmd.AddCommand(&cobra.Command{
		Use:   "delete <key>",
		Short a default entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := defaults.NewStore			if err != nil {
				return err
			}
			return s.Delete(args[0])
		},
	})

	defaultsCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all default values",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := defaults.NewStore(defaultsStorePath)
			if err != nil {
				return err
			}
			all := s.All()
			if len(all) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no defaults set")
				return nil
			}
			keys := make([]string, 0, len(all))
			for k := range all {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, all[k])
			}
			return nil
		},
	})

	_ = os.MkdirAll(defaultsStorePath, 0o755)
	rootCmd.AddCommand(defaultsCmd)
}
