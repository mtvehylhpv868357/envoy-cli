// Package validate provides environment variable validation against a schema.
package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Key      string `json:"key"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
	MinLen   int    `json:"min_len,omitempty"`
	MaxLen   int    `json:"max_len,omitempty"`
}

// Issue represents a single validation problem.
type Issue struct {
	Key     string
	Message string
	Level   string // "error" or "warn"
}

// Schema holds a collection of rules.
type Schema struct {
	Rules []Rule `json:"rules"`
}

// Validate checks the provided env map against the schema rules.
func Validate(env map[string]string, schema Schema) []Issue {
	var issues []Issue

	for _, rule := range schema.Rules {
		val, exists := env[rule.Key]

		if rule.Required && !exists {
			issues = append(issues, Issue{
				Key:     rule.Key,
				Message: "required variable is missing",
				Level:   "error",
			})
			continue
		}

		if !exists {
			continue
		}

		if rule.MinLen > 0 && len(val) < rule.MinLen {
			issues = append(issues, Issue{
				Key:     rule.Key,
				Message: fmt.Sprintf("value too short (min %d chars)", rule.MinLen),
				Level:   "error",
			})
		}

		if rule.MaxLen > 0 && len(val) > rule.MaxLen {
			issues = append(issues, Issue{
				Key:     rule.Key,
				Message: fmt.Sprintf("value too long (max %d chars)", rule.MaxLen),
				Level:   "warn",
			})
		}

		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				issues = append(issues, Issue{
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
					Level:   "error",
				})
				continue
			}
			if !re.MatchString(val) {
				issues = append(issues, Issue{
					Key:     rule.Key,
					Message: fmt.Sprintf("value does not match pattern %q", rule.Pattern),
					Level:   "error",
				})
			}
		}
	}

	return issues
}

// HasErrors returns true if any issue has level "error".
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if strings.EqualFold(i.Level, "error") {
			return true
		}
	}
	return false
}
