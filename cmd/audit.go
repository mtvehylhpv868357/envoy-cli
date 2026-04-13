package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envoy-cli/internal/audit"
	"envoy-cli/internal/config"
)

func init() {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "View the audit log of envoy-cli actions",
	}

	auditListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all recorded audit events",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			logger, err := audit.NewLogger(cfg.DataDir)
			if err != nil {
				return fmt.Errorf("open audit log: %w", err)
			}
			entries, err := logger.ReadAll()
			if err != nil {
				return fmt.Errorf("read audit log: %w", err)
			}
			if len(entries) == 0 {
				fmt.Println("No audit events recorded.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIMESTAMP\tEVENT\tPROJECT\tDETAIL")
			for _, e := range entries {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					e.Timestamp.Format("2006-01-02 15:04:05"),
					e.Event, e.Project, e.Detail)
			}
			return w.Flush()
		},
	}

	auditClearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all audit log entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}
			path := cfg.DataDir + "/audit.log"
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("clear audit log: %w", err)
			}
			fmt.Println("Audit log cleared.")
			return nil
		},
	}

	auditCmd.AddCommand(auditListCmd, auditClearCmd)
	rootCmd.AddCommand(auditCmd)
}
