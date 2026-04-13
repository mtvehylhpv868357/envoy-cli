package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envoy-cli/internal/profile"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage environment variable profiles",
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := profile.LoadStore("")
		if err != nil {
			return fmt.Errorf("failed to load store: %w", err)
		}
		profiles := store.List()
		if len(profiles) == 0 {
			fmt.Println("No profiles found. Use 'envoy profile add' to create one.")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tACTIVE")
		for _, p := range profiles {
			active := ""
			if p == store.Active() {
				active = "*"
			}
			fmt.Fprintf(w, "%s\t%s\n", p, active)
		}
		return w.Flush()
	},
}

var profileUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Set the active profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := profile.LoadStore("")
		if err != nil {
			return fmt.Errorf("failed to load store: %w", err)
		}
		if err := store.SetActive(args[0]); err != nil {
			return fmt.Errorf("failed to set active profile: %w", err)
		}
		fmt.Printf("Switched to profile '%s'\n", args[0])
		return nil
	},
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := profile.LoadStore("")
		if err != nil {
			return fmt.Errorf("failed to load store: %w", err)
		}
		if err := store.Delete(args[0]); err != nil {
			return fmt.Errorf("failed to delete profile: %w", err)
		}
		fmt.Printf("Deleted profile '%s'\n", args[0])
		return nil
	},
}

func init() {
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileUseCmd)
	profileCmd.AddCommand(profileDeleteCmd)
	rootCmd.AddCommand(profileCmd)
}
