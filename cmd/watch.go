package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/envoy-cli/internal/watch"
)

func init() {
	var interval int

	watchCmd := &cobra.Command{
		Use:   "watch [file]",
		Short: "Watch a .env file and print a notice on change",
		Long: `Polls the specified .env file (or the active profile file) for
modifications and prints a notification each time a change is detected.

Press Ctrl+C to stop watching.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var target string
			if len(args) == 1 {
				target = args[0]
			} else {
				target = ".env"
			}

			abs, err := filepath.Abs(target)
			if err != nil {
				return fmt.Errorf("resolving path: %w", err)
			}
			if _, err := os.Stat(abs); err != nil {
				return fmt.Errorf("file not found: %s", abs)
			}

			duration := time.Duration(interval) * time.Millisecond
			w := watch.New(duration, func(e watch.Event) {
				fmt.Fprintf(cmd.OutOrStdout(), "[envoy] change detected: %s (modified %s)\n",
					e.Path, e.ModTime.Format("15:04:05"))
			})

			if err := w.Add(abs); err != nil {
				return fmt.Errorf("watch: %w", err)
			}
			w.Start()
			defer w.Stop()

			fmt.Fprintf(cmd.OutOrStdout(), "Watching %s (every %dms). Press Ctrl+C to stop.\n", abs, interval)

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			<-sig
			fmt.Fprintln(cmd.OutOrStdout(), "\nStopped.")
			return nil
		},
	}

	watchCmd.Flags().IntVarP(&interval, "interval", "i", 500, "poll interval in milliseconds")
	rootCmd.AddCommand(watchCmd)
}
