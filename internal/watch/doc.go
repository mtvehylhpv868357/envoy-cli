// Package watch implements lightweight poll-based file watching for envoy-cli.
//
// It is used to automatically detect changes to .env profile files so the
// CLI can notify the user or trigger a reload without requiring inotify or
// OS-specific APIs.
//
// Usage:
//
//	w := watch.New(500*time.Millisecond, func(e watch.Event) {
//		fmt.Println("changed:", e.Path)
//	})
//	_ = w.Add(".env")
//	w.Start()
//	defer w.Stop()
package watch
