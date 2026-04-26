// Package envcast provides utilities for casting environment variable
// string values to typed Go values such as int, bool, float64, and duration.
package envcast

import (
	"fmt"
	"strconv"
	"time"
)

// Result holds a cast value along with any error that occurred.
type Result struct {
	Raw   string
	Error error
}

// ToString returns the raw string value unchanged.
func ToString(v string) (string, error) {
	return v, nil
}

// ToBool parses "true", "1", "yes" as true and "false", "0", "no" as false.
func ToBool(v string) (bool, error) {
	switch v {
	case "true", "1", "yes", "TRUE", "YES":
		return true, nil
	case "false", "0", "no", "FALSE", "NO":
		return false, nil
	}
	return false, fmt.Errorf("envcast: cannot convert %q to bool", v)
}

// ToInt parses the string as a base-10 integer.
func ToInt(v string) (int, error) {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("envcast: cannot convert %q to int: %w", v, err)
	}
	return n, nil
}

// ToFloat64 parses the string as a 64-bit floating-point number.
func ToFloat64(v string) (float64, error) {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, fmt.Errorf("envcast: cannot convert %q to float64: %w", v, err)
	}
	return f, nil
}

// ToDuration parses the string as a time.Duration (e.g. "5s", "2m", "1h").
func ToDuration(v string) (time.Duration, error) {
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("envcast: cannot convert %q to duration: %w", v, err)
	}
	return d, nil
}

// Map applies a cast function to every value in the provided map, returning
// a new map of the same keys with converted values. Keys whose values fail
// conversion are collected into a combined error.
func Map[T any](vars map[string]string, fn func(string) (T, error)) (map[string]T, error) {
	out := make(map[string]T, len(vars))
	var errs []string
	for k, v := range vars {
		cast, err := fn(v)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", k, err))
			continue
		}
		out[k] = cast
	}
	if len(errs) > 0 {
		return out, fmt.Errorf("envcast: %d conversion error(s): %v", len(errs), errs)
	}
	return out, nil
}
