// Package audit provides append-only structured logging for envoy-cli actions.
//
// Each operation that mutates state (profile switches, vault encryption,
// snapshot saves, etc.) should record an Entry via Logger.Log so that
// users can inspect a full history of changes with `envoy audit list`.
//
// Log files are stored as newline-delimited JSON in the configured data
// directory (default: ~/.config/envoy-cli/) with permissions 0600.
//
// # Entry Format
//
// Each log entry is a JSON object written as a single line and contains
// at minimum the following fields:
//
//   - timestamp: RFC 3339 UTC time at which the event occurred
//   - action:    short identifier for the operation (e.g. "profile.switch")
//   - actor:     OS user that triggered the action
//   - details:   map of action-specific key/value pairs
//
// # Retention
//
// Log rotation and pruning are not performed automatically. Users are
// responsible for managing log size. A future `envoy audit prune` command
// may be added to remove entries older than a configurable duration.
package audit
