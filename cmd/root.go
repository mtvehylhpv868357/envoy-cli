package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envoy",
	Short: "envoy — manage and switch environment variable profiles",
	Long: `envoy is a lightweight CLI for managing environment variable profiles.

Use profiles to store, switch, and export environment variables across projects.
Run 'envoy profile list' to see available profiles, or 'envoy shell' to emit
export statements for your current shell session.`,
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
