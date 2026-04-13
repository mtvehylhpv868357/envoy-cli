package template

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// varPattern matches ${VAR_NAME} or $VAR_NAME style placeholders
var varPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Result holds the output of a template render operation
type Result struct {
	Rendered string
	Missing  []string
}

// Render replaces environment variable placeholders in the template string
// using the provided vars map. Missing variables are collected in Result.Missing.
func Render(tmpl string, vars map[string]string) Result {
	missing := []string{}
	seen := map[string]bool{}

	rendered := varPattern.ReplaceAllStringFunc(tmpl, func(match string) string {
		name := extractName(match)
		if val, ok := vars[name]; ok {
			return val
		}
		if !seen[name] {
			missing = append(missing, name)
			seen[name] = true
		}
		return match
	})

	return Result{Rendered: rendered, Missing: missing}
}

// RenderFromEnv renders a template using the current process environment.
func RenderFromEnv(tmpl string) Result {
	vars := map[string]string{}
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return Render(tmpl, vars)
}

// extractName pulls the variable name from a matched placeholder string.
func extractName(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}

// ValidatePlaceholders returns an error if any placeholders in the template
// are malformed (e.g. lowercase names).
func ValidatePlaceholders(tmpl string) error {
	badPattern := regexp.MustCompile(`\$\{([^}]*)\}|\$([a-z][a-zA-Z0-9_]*)`)
	matches := badPattern.FindAllString(tmpl, -1)
	if len(matches) > 0 {
		return fmt.Errorf("template contains non-uppercase placeholders: %v", matches)
	}
	return nil
}
