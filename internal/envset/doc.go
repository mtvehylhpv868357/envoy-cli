// Package envset provides set-theoretic operations on collections of
// environment variable maps: union, intersection, difference, and
// symmetric difference.
//
// These operations are useful when combining profiles, computing what
// variables are shared across environments, or finding variables unique
// to a specific profile.
//
// Example:
//
//	base := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
//	prod := map[string]string{"DB_HOST": "prod.db", "API_KEY": "secret"}
//
//	// Keys only in base:
//	onlyBase := envset.Difference(envset.DefaultOptions(), base, prod)
//
//	// Keys in both:
//	shared := envset.Intersect(envset.DefaultOptions(), base, prod)
package envset
