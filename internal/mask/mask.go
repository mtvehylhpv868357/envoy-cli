// Package mask provides utilities for masking sensitive environment variable values
// before display or logging.
package mask

import "strings"

// DefaultOptions returns a MaskOptions with sensible defaults.
func DefaultOptions() Options {
	return Options{
		RevealChars: 4,
		MaskChar:    '*',
		SensitiveKeys: []string{
			"SECRET", "PASSWORD", "PASSWD", "TOKEN", "KEY",
			"PRIVATE", "CREDENTIAL", "AUTH", "APIKEY", "API_KEY",
		},
	}
}

// Options configures masking behaviour.
type Options struct {
	// RevealChars is the number of trailing characters to reveal.
	RevealChars int
	// MaskChar is the character used to replace hidden characters.
	MaskChar rune
	// SensitiveKeys contains substrings that mark a key as sensitive.
	SensitiveKeys []string
}

// IsSensitive reports whether the given key should be treated as sensitive
// based on the configured SensitiveKeys list (case-insensitive).
func (o Options) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, s := range o.SensitiveKeys {
		if strings.Contains(upper, strings.ToUpper(s)) {
			return true
		}
	}
	return false
}

// Value masks a single value, optionally revealing the last RevealChars characters.
func (o Options) Value(v string) string {
	if len(v) == 0 {
		return ""
	}
	reveal := o.RevealChars
	if reveal >= len(v) {
		reveal = 0
	}
	hidden := len(v) - reveal
	return strings.Repeat(string(o.MaskChar), hidden) + v[hidden:]
}

// Vars returns a copy of the vars map with sensitive values masked.
func (o Options) Vars(vars map[string]string) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		if o.IsSensitive(k) {
			out[k] = o.Value(v)
		} else {
			out[k] = v
		}
	}
	return out
}
