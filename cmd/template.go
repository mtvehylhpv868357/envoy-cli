package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envoy-cli/internal/env"
	"github.com/yourorg/envoy-cli/internal/profile"
	"github.com/yourorg/envoy-cli/internal/template"
)

func init() {
	var profileName string
	var strict bool

	templateCmd := &cobra.Command{
		Use:   "template <file>",
		Short: "Render a template file using a profile's environment variables",
		Long: `Reads a template file and substitutes $VAR or ${VAR} placeholders
using the variables from the specified (or active) profile.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("reading template file: %w", err)
			}

			store, err := profile.LoadStore("")
			if err != nil {
				return fmt.Errorf("loading profile store: %w", err)
			}

			name := profileName
			if name == "" {
				name = store.Active
			}
			if name == "" {
				return fmt.Errorf("no profile specified and no active profile set")
			}

			p, err := store.Get(name)
			if err != nil {
				return fmt.Errorf("profile %q not found: %w", name, err)
			}

			vars := env.Apply(p.Vars)
			result := template.Render(string(data), vars)

			if strict && len(result.Missing) > 0 {
				return fmt.Errorf("unresolved placeholders: %s", strings.Join(result.Missing, ", "))
			}

			if len(result.Missing) > 0 {
				fmt.Fprintf(os.Stderr, "warning: unresolved placeholders: %s\n", strings.Join(result.Missing, ", "))
			}

			fmt.Print(result.Rendered)
			return nil
		},
	}

	templateCmd.Flags().StringVarP(&profileName, "profile", "p", "", "Profile to use for variable substitution")
	templateCmd.Flags().BoolVar(&strict, "strict", false, "Fail if any placeholders are unresolved")

	rootCmd.AddCommand(templateCmd)
}
