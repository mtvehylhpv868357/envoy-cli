// Package lint provides validation for environment variable profiles.
package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// Issue represents a single lint warning or error found in a profile.
type Issue struct {
	Key      string
	Severity string // "error" or "warn"
	Message  string
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(i.Severity), i.Key, i.Message)
}

var validKeyRe = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// Check validates a map of environment variables and returns a list of issues.
func Check(vars map[string]string) []Issue {
	var issues []Issue

	for k, v := range vars {
		// Key naming convention
		if !validKeyRe.MatchString(k) {
			issues = append(issues, Issue{
				Key:      k,
				Severity: "warn",
				Message:  "key should be UPPER_SNAKE_CASE",
			})
		}

		// Empty value
		if strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{
				Key:      k,
				Severity: "warn",
				Message:  "value is empty",
			})
		}

		// Potential secret stored in plaintext
		lower := strings.ToLower(k)
		if containsAny(lower, []string{"secret", "password", "passwd", "token", "apikey", "api_key"}) {
			issues = append(issues, Issue{
				Key:      k,
				Severity: "error",
				Message:  "possible secret stored in plaintext — consider using vault encrypt",
			})
		}
	}

	return issues
}

// HasErrors returns true if any issue has severity "error".
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if i.Severity == "error" {
			return true
		}
	}
	return false
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
