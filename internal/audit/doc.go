// Package audit provides append-only structured logging for envoy-cli actions.
//
// Each operation that mutates state (profile switches, vault encryption,
// snapshot saves, etc.) should record an Entry via Logger.Log so that
// users can inspect a full history of changes with `envoy audit list`.
//
// Log files are stored as newline-delimited JSON in the configured data
// directory (default: ~/.config/envoy-cli/) with permissions 0600.
package audit
