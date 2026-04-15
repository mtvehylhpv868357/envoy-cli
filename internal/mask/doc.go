// Package mask provides utilities for detecting and masking sensitive
// environment variable values before they are displayed to the user or
// written to log output.
//
// A key is considered sensitive when its name contains one of the
// configured substrings (e.g. PASSWORD, TOKEN, SECRET).  When masking,
// the last N characters of the value are optionally revealed so the user
// can still distinguish between different values without exposing the
// full secret.
//
// Example:
//
//	opts := mask.DefaultOptions()
//	masked := opts.Vars(profile.Vars)
package mask
