package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/profile"
	"github.com/yourorg/envoy-cli/internal/validate"
)

func init() {
	var schemaPath string
	var storePath string

	validateCmd := &cobra.Command{
		Use:   "validate [profile]",
		Short: "Validate a profile's variables against a schema",
		Long: `Validate checks the active (or specified) profile's environment variables
against rules defined in a schema file (default: .envoy-schema.json).`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := profile.LoadStore(storePath)
			if err != nil {
				return fmt.Errorf("load store: %w", err)
			}

			var name string
			if len(args) == 1 {
				name = args[0]
			} else {
				name = store.Active
			}
			if name == "" {
				return fmt.Errorf("no profile specified and no active profile set")
			}

			p, err := store.Get(name)
			if err != nil {
				return fmt.Errorf("get profile %q: %w", name, err)
			}

			if schemaPath == "" {
				schemaPath = validate.DefaultSchemaPath(".")
			}
			schema, err := validate.LoadSchema(schemaPath)
			if err != nil {
				return fmt.Errorf("load schema: %w", err)
			}

			issues := validate.Validate(p.Vars, schema)
			if len(issues) == 0 {
				fmt.Printf("✓ Profile %q passed all validation rules.\n", name)
				return nil
			}

			for _, iss := range issues {
				icon := "⚠"
				if iss.Level == "error" {
					icon = "✗"
				}
				fmt.Fprintf(os.Stderr, "%s [%s] %s: %s\n", icon, iss.Level, iss.Key, iss.Message)
			}

			if validate.HasErrors(issues) {
				return fmt.Errorf("validation failed with errors")
			}
			return nil
		},
	}

	validateCmd.Flags().StringVar(&schemaPath, "schema", "", "path to schema file (default: .envoy-schema.json)")
	validateCmd.Flags().StringVar(&storePath, "store", "", "path to profile store directory")

	rootCmd.AddCommand(validateCmd)
}
