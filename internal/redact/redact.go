// Package redact provides utilities for redacting sensitive environment
// variable values before displaying or logging them.
package redact

import (
	"strings"
)

// DefaultOptions returns a RedactOptions with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Replacement: "[REDACTED]",
		PartialReveal: false,
		RevealChars: 4,
	}
}

// Options controls how values are redacted.
type Options struct {
	// Replacement is the string used to replace sensitive values.
	Replacement string
	// PartialReveal shows the last N characters of the value.
	PartialReveal bool
	// RevealChars is the number of trailing characters to reveal when PartialReveal is true.
	RevealChars int
}

// sensitiveKeys contains substrings that indicate a key holds sensitive data.
var sensitiveKeys = []string{
	"secret", "password", "passwd", "token", "apikey", "api_key",
	"private", "credential", "auth", "cert", "key",
}

// IsSensitive returns true if the key name suggests a sensitive value.
func IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}

// Value redacts a single value according to the provided options.
// If the value is short (<= RevealChars), it is fully replaced.
func Value(v string, opts Options) string {
	if !opts.PartialReveal || len(v) <= opts.RevealChars {
		return opts.Replacement
	}
	return opts.Replacement + "..." + v[len(v)-opts.RevealChars:]
}

// Map returns a copy of the input map with sensitive values redacted.
func Map(vars map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if IsSensitive(k) {
			out[k] = Value(v, opts)
		} else {
			out[k] = v
		}
	}
	return out
}
