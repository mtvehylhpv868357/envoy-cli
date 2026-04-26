package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/envoy-cli/internal/envcast"
	"github.com/your-org/envoy-cli/internal/profile"
)

func init() {
	castCmd := &cobra.Command{
		Use:   "cast <profile> <type>",
		Short: "Display profile variables cast to a given type",
		Long: `Cast all (or selected) environment variables in a profile to the
specified type and print the results.

Supported types: string, bool, int, float, duration`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			profName := args[0]
			castType := strings.ToLower(args[1])

			storeDir, _ := cmd.Flags().GetString("store")
			if storeDir == "" {
				storeDir = defaultProfileDir()
			}

			st, err := profile.LoadStore(storeDir)
			if err != nil {
				return fmt.Errorf("load store: %w", err)
			}
			vars, err := st.Get(profName)
			if err != nil {
				return fmt.Errorf("profile %q not found: %w", profName, err)
			}

			keys := make([]string, 0, len(vars))
			for k := range vars {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				raw := vars[k]
				result, castErr := castValue(raw, castType)
				if castErr != nil {
					fmt.Fprintf(os.Stderr, "  %s: error: %v\n", k, castErr)
					continue
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %s=%v\n", k, result)
			}
			return nil
		},
	}

	castCmd.Flags().String("store", "", "path to profile store directory")
	rootCmd.AddCommand(castCmd)
}

func castValue(raw, castType string) (interface{}, error) {
	switch castType {
	case "string":
		return envcast.ToString(raw)
	case "bool":
		return envcast.ToBool(raw)
	case "int":
		return envcast.ToInt(raw)
	case "float", "float64":
		return envcast.ToFloat64(raw)
	case "duration":
		return envcast.ToDuration(raw)
	default:
		return nil, fmt.Errorf("unknown type %q; supported: string, bool, int, float, duration", castType)
	}
}
