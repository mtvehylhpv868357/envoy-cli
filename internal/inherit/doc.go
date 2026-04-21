// Package inherit implements profile inheritance for envoy-cli.
//
// A child profile is created by copying all variables from a base profile and
// then applying a caller-supplied overrides map on top.  This allows teams to
// maintain a canonical "base" environment and derive per-developer or
// per-environment variants without duplicating every variable.
//
// Example usage:
//
//	opts := inherit.DefaultOptions()
//	opts.Overwrite = true
//	err := inherit.Profile(store, "production", "staging", map[string]string{
//		"API_URL": "https://staging.example.com",
//	}, opts)
package inherit
