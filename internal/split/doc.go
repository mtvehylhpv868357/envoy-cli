// Package split provides the ability to decompose a single large environment
// profile into several smaller, focused profiles by matching key prefixes.
//
// Example usage:
//
//	opts := split.DefaultOptions()
//	created, err := split.ByPrefix(store, "monolith", []string{"DB_", "APP_", "CACHE_"}, opts)
//
// Each prefix becomes its own profile named after the lower-cased prefix
// (with a trailing underscore stripped). Keys are optionally stripped of the
// matched prefix so the resulting profiles contain clean, context-free names.
package split
