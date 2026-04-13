package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/profile"
	"envoy-cli/internal/run"
)

func init() {
	var profileName string
	var noInherit bool

	runCmd := &cobra.Command{
		Use:   "run [flags] -- <command> [args...]",
		Short: "Run a command with an injected environment profile",
		Long: `Run executes a command with the variables from the specified
profile (or the active profile) injected into its environment.

Example:
  envoy run -- node server.js
  envoy run --profile staging -- make deploy`,
		DisableFlagParsing: false,
		Args:               cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := profile.LoadStore("")
			if err != nil {
				return fmt.Errorf("failed to load profile store: %w", err)
			}

			name := profileName
			if name == "" {
				name = store.Active
			}
			if name == "" {
				return fmt.Errorf("no profile specified and no active profile set; use --profile or `envoy profile use <name>`")
			}

			p, err := store.Get(name)
			if err != nil {
				return fmt.Errorf("profile %q not found: %w", name, err)
			}

			opts := run.DefaultOptions()
			opts.Vars = p.Vars
			opts.Inherit = !noInherit

			if err := run.Run(args, opts); err != nil {
				// Preserve exit code when possible
				fmt.Fprintf(os.Stderr, "envoy run: %v\n", err)
				os.Exit(1)
			}
			return nil
		},
	}

	runCmd.Flags().StringVarP(&profileName, "profile", "p", "", "Profile name to inject (defaults to active profile)")
	runCmd.Flags().BoolVar(&noInherit, "no-inherit", false, "Do not inherit the current process environment")

	rootCmd.AddCommand(runCmd)
}
