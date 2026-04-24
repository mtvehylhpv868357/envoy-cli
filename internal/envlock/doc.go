// Package envlock implements profile locking for envoy-cli.
//
// A locked profile cannot be modified by commands that respect lock state.
// Lock entries are stored as JSON files in a dedicated directory and record
// the profile name, the time the lock was applied, and an optional reason.
//
// Typical usage:
//
//	s, err := envlock.NewStore("/path/to/locks")
//	if err != nil { ... }
//
//	// Lock a profile before a production deploy.
//	s.Lock("production", "deploy freeze")
//
//	// Check before mutating.
//	if s.IsLocked("production") {
//		fmt.Println("profile is locked")
//	}
//
//	// Release the lock afterwards.
//	s.Unlock("production")
package envlock
