// Package archive provides Pack and Unpack helpers for bundling environment
// profiles into portable gzipped tar archives.
//
// Use Pack to export one or more profiles from a store directory into a
// writer, and Unpack to restore them into a destination directory.
//
// Example:
//
//	var buf bytes.Buffer
//	archive.Pack("/home/user/.envoy/profiles", []string{"dev", "staging"}, &buf)
//	// write buf to a file or transfer it
//	archive.Unpack(&buf, "/home/user/.envoy/profiles", false)
package archive
