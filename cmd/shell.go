package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/env"
	"envoy-cli/internal/profile"
	"envoy-cli/internal/shell"
)

var shellTypeFlag string

var shellCmd = &cobra.Command{
	Use:   "shell [profile]",
	Short: "Emit shell export commands for a profile's variables",
	Long: `Prints shell-specific export statements for the given profile.
Eval the output to apply variables in the current shell session:

  eval $(envoy shell myprofile)
`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := profile.LoadStore("")
		if err != nil {
			return fmt.Errorf("loading store: %w", err)
		}

		profileName := store.Active
		if len(args) == 1 {
			profileName = args[0]
		}
		if profileName == "" {
			return fmt.Errorf("no profile specified and no active profile set")
		}

		p, err := store.Get(profileName)
		if err != nil {
			return fmt.Errorf("profile %q not found: %w", profileName, err)
		}

		vars, err := env.ParseDotEnv(p.RawEnv)
		if err != nil {
			return fmt.Errorf("parsing env: %w", err)
		}

		sh := shell.Detect()
		if shellTypeFlag != "" {
			sh = shell.ShellType(shellTypeFlag)
		}

		fmt.Fprint(os.Stdout, shell.ExportScript(vars, sh))
		return nil
	},
}

func init() {
	shellCmd.Flags().StringVar(&shellTypeFlag, "shell", "", "override shell type (bash, zsh, fish)")
	rootCmd.AddCommand(shellCmd)
}
