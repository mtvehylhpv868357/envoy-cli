// Package expire provides TTL-based expiration tracking for environment profiles.
//
// It allows setting an expiration time on a named profile so that downstream
// tooling can warn or refuse to apply variables from a profile that has passed
// its intended lifetime.
//
// Usage:
//
//	store, err := expire.NewStore("/path/to/expire-store")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Set a profile to expire in 24 hours
//	err = store.Set("production", time.Now().Add(24*time.Hour))
//
//	// Check if a profile is expired
//	expired, err := store.IsExpired("production")
//
//	// List all tracked profiles with their expiry times
//	entries, err := store.List()
package expire
