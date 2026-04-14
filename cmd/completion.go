package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/completion"
)

func init() {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for envoy-cli.

Source the output to enable tab-completion in your shell:

  # bash
  source <(envoy completion bash)

  # zsh
  source <(envoy completion zsh)

  # fish
  envoy completion fish | source`,
	}

	completionCmd.AddCommand(&cobra.Command{
		Use:   "bash",
		Short: "Generate bash completion script",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenBashCompletion(os.Stdout)
		},
	})

	completionCmd.AddCommand(&cobra.Command{
		Use:   "zsh",
		Short: "Generate zsh completion script",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenZshCompletion(os.Stdout)
		},
	})

	completionCmd.AddCommand(&cobra.Command{
		Use:   "fish",
		Short: "Generate fish completion script",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		},
	})

	completionCmd.AddCommand(&cobra.Command{
		Use:   "profiles",
		Short: "List profile names for scripted completion",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			storeDir := filepath.Join(home, ".config", "envoy", "profiles")
			prefix := ""
			if len(args) > 0 {
				prefix = args[0]
			}
			names := completion.FilterPrefix(completion.ProfileNames(storeDir), prefix)
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
			return nil
		},
	})

	rootCmd.AddCommand(completionCmd)
}
