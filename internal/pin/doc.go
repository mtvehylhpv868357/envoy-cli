// Package pin manages directory-to-profile pins for envoy-cli.
//
// A "pin" associates an absolute directory path with a named environment
// profile. When envoy-cli is invoked from a pinned directory (or one of its
// children), it can automatically activate the associated profile without
// requiring the user to specify one explicitly.
//
// Pins are stored as a JSON file (default: ~/.config/envoy/pins.json) and
// manipulated via the `envoy pin` sub-commands:
//
//	envoy pin set   <profile>          # pin cwd to profile
//	envoy pin get                      # show profile pinned to cwd
//	envoy pin remove                   # remove pin for cwd
//	envoy pin list                     # list all pins
package pin
