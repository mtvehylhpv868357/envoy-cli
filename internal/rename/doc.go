// Package rename implements profile renaming for envoy-cli.
//
// It provides a single entry-point, Profile, which atomically renames an
// existing environment profile within a store.  The operation preserves the
// active-profile pointer when the renamed profile was previously active.
//
// Example usage:
//
//	err := rename.Profile(store, "dev", "development", rename.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
package rename
