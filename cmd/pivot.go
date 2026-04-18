package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/envoy-cli/internal/pivot"
	"github.com/envoy-cli/internal/profile"
)

func init() {
	var storeDir string
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "pivot [profile1] [profile2] ...",
		Short: "Transpose profiles into a key-centric comparison table",
		Long:  "Show all keys across the given profiles side by side, highlighting where values differ or are missing.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := profile.LoadStore(storeDir)
			if err != nil {
				return fmt.Errorf("load store: %w", err)
			}

			profiles := map[string]map[string]string{}
			for _, name := range args {
				p, err := store.Get(name)
				if err != nil {
					return fmt.Errorf("profile %q not found: %w", name, err)
				}
				profiles[name] = p.Vars
			}

			rows := pivot.Profiles(profiles, pivot.DefaultOptions())

			if outputJSON {
				return json.NewEncoder(os.Stdout).Encode(rows)
			}

			// Header
			fmt.Printf("%-30s", "KEY")
			for _, name := range args {
				fmt.Printf("  %-20s", name)
			}
			fmt.Println()
			fmt.Println(strings.Repeat("-", 30+len(args)*22))

			for _, row := range rows {
				fmt.Printf("%-30s", row.Key)
				for _, name := range args {
					v, ok := row.Values[name]
					if !ok {
						v = "<missing>"
					}
					if len(v) > 18 {
						v = v[:15] + "..."
					}
					fmt.Printf("  %-20s", v)
				}
				fmt.Println()
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&storeDir, "store", "", "profile store directory")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "output as JSON")
	rootCmd.AddCommand(cmd)
}
