// Package search implements cross-profile environment variable search.
//
// It allows users to find which profiles contain a given key or value,
// supporting substring matching and exact-key lookups.
//
// Example usage:
//
//	store, _ := profile.LoadStore(dir)
//	results, err := search.Profiles(store, search.Options{
//		KeyPattern: "DB_",
//	})
//	for _, r := range results {
//		fmt.Printf("%s\t%s=%s\n", r.Profile, r.Key, r.Value)
//	}
package search
