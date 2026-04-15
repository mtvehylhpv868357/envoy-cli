// Package doctor provides environment profile health checks.
package doctor

import (
	"fmt"
	"os"
	"strings"
)

// Severity indicates the level of a diagnostic finding.
type Severity string

const (
	SeverityOK      Severity = "ok"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Finding represents a single diagnostic result.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

// Report holds all findings from a diagnostic run.
type Report struct {
	Findings []Finding
}

// HasErrors returns true if any finding has error severity.
func (r *Report) HasErrors() bool {
	for _, f := range r.Findings {
		if f.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary line.
func (r *Report) Summary() string {
	var errs, warns int
	for _, f := range r.Findings {
		switch f.Severity {
		case SeverityError:
			errs++
		case SeverityWarning:
			warn++
		}
	}
	return fmt.Sprintf("%d error(s), %d warning(s)", errs, warns)
}

// Check runs all health checks against the provided env vars map.
func Check(vars map[string]string) *Report {
	report := &Report{}

	for k, v := range vars {
		// Check for empty values
		if strings.TrimSpace(v) == "" {
			report.Findings = append(report.Findings, Finding{
				Key:      k,
				Message:  "value is empty",
				Severity: SeverityWarning,
			})
		}

		// Check for keys referencing unset OS vars
		if strings.HasPrefix(v, "$") {
			ref := strings.TrimPrefix(v, "$")
			if os.Getenv(ref) == "" {
				report.Findings = append(report.Findings, Finding{
					Key:      k,
					Message:  fmt.Sprintf("references unset variable $%s", ref),
					Severity: SeverityWarning,
				})
			}
		}

		// Check for lowercase keys (convention warning)
		if k != strings.ToUpper(k) {
			report.Findings = append(report.Findings, Finding{
				Key:      k,
				Message:  "key is not uppercase (convention)",
				Severity: SeverityWarning,
			})
		}

		// Check for keys that look like secrets but may be plaintext
		upper := strings.ToUpper(k)
		if containsSensitive(upper) && len(v) < 8 {
			report.Findings = append(report.Findings, Finding{
				Key:      k,
				Message:  "looks like a secret but value is suspiciously short",
				Severity: SeverityError,
			})
		}
	}

	return report
}

func containsSensitive(key string) bool {
	sensitive := []string{"SECRET", "PASSWORD", "TOKEN", "PRIVATE_KEY", "API_KEY"}
	for _, s := range sensitive {
		if strings.Contains(key, s) {
			return true
		}
	}
	return false
}
