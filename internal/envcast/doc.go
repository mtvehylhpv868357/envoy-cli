// Package envcast provides type-safe casting utilities for environment
// variable string values.
//
// Supported target types:
//
//   - string  (passthrough)
//   - bool    ("true"/"false", "1"/"0", "yes"/"no")
//   - int     (base-10 integer)
//   - float64 (IEEE 754 double precision)
//   - time.Duration (Go duration string, e.g. "5s", "2m")
//
// The Map helper applies any cast function to an entire map[string]string,
// collecting conversion errors without halting early.
package envcast
