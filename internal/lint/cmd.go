package lint

// Severity levels for lint issues.
const (
	SeverityWarn  = "WARN"
	SeverityError = "ERROR"
)

// Issue represents a single lint finding for an environment variable.
type Issue struct {
	Key      string
	Message  string
	Severity string
}

// Report holds the full result of a lint run.
type Report struct {
	Issues []Issue
}

// Add appends an issue to the report.
func (r *Report) Add(key, message, severity string) {
	r.Issues = append(r.Issues, Issue{Key: key, Message: message, Severity: severity})
}

// HasErrors returns true if any issue has ERROR severity.
func (r *Report) HasErrors() bool {
	for _, i := range r.Issues {
		if i.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary string.
func (r *Report) Summary() string {
	if len(r.Issues) == 0 {
		return "No lint issues found."
	}
	errCount, warnCount := 0, 0
	for _, i := range r.Issues {
		if i.Severity == SeverityError {
			errCount++
		} else {
			warnCount++
		}
	}
	return fmt.Sprintf("%d error(s), %d warning(s) found.", errCount, warnCount)
}
