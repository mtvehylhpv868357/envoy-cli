package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/envlock"
)

var lockStorePath = filepath.Join(os.Getenv("HOME"), ".config", "envoy", "locks")

func init() {
	lockCmd := &cobra.Command{
		Use:   "lock",
		Short: "Lock or unlock environment profiles to prevent modifications",
	}

	var reason string

	setCmd := &cobra.Command{
		Use:   "set <profile>",
		Short: "Lock a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := envlock.NewStore(lockStorePath)
			if err != nil {
				return err
			}
			if err := s.Lock(args[0], reason); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "profile %q locked\n", args[0])
			return nil
		},
	}
	setCmd.Flags().StringVarP(&reason, "reason", "r", "", "optional reason for locking")

	unsetCmd := &cobra.Command{
		Use:   "unset <profile>",
		Short: "Unlock a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := envlock.NewStore(lockStorePath)
			if err != nil {
				return err
			}
			if err := s.Unlock(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "profile %q unlocked\n", args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all locked profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := envlock.NewStore(lockStorePath)
			if err != nil {
				return err
			}
			locks, err := s.List()
			if err != nil {
				return err
			}
			if len(locks) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no locked profiles")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "PROFILE\tLOCKED AT\tREASON")
			for _, l := range locks {
				fmt.Fprintf(w, "%s\t%s\t%s\n", l.Profile, l.LockedAt.Format("2006-01-02 15:04:05"), l.Reason)
			}
			return w.Flush()
		},
	}

	lockCmd.AddCommand(setCmd, unsetCmd, listCmd)
	rootCmd.AddCommand(lockCmd)
}
