package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/envoy-cli/envoy-cli/internal/history"
	"github.com/spf13/cobra"
)

func historyStorePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envoy", "history.json")
}

func init() {
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show profile activation history",
		Long:  "Display a chronological log of profile activations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := history.NewStore(historyStorePath())
			if err != nil {
				return fmt.Errorf("opening history store: %w", err)
			}
			entries, err := s.ReadAll()
			if err != nil {
				return fmt.Errorf("reading history: %w", err)
			}
			if len(entries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No history recorded yet.")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "#\tPROFILE\tACTIVATED")
			for i, e := range entries {
				fmt.Fprintf(w, "%d\t%s\t%s\n", i+1, e.Profile, e.Activated.Local().Format(time.RFC822))
			}
			return w.Flush()
		},
	}

	historyClearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all history entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := history.NewStore(historyStorePath())
			if err != nil {
				return fmt.Errorf("opening history store: %w", err)
			}
			if err := s.Clear(); err != nil {
				return fmt.Errorf("clearing history: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "History cleared.")
			return nil
		},
	}

	historyCmd.AddCommand(historyClearCmd)
	rootCmd.AddCommand(historyCmd)
}
