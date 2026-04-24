// Package stats computes aggregated statistics over a collection of
// environment variable profiles.
//
// Given a map of profile names to their key/value pairs, Compute returns
// a Report containing:
//
//   - ProfileCount  – number of profiles analysed
//   - TotalVars     – sum of all key/value pairs across profiles
//   - UniqueKeys    – keys that appear in exactly one profile
//   - SharedKeys    – keys that appear in more than one profile
//   - EmptyValues   – count of blank values across all profiles
//   - SensitiveKeys – distinct keys whose names suggest sensitive data
//   - KeyFrequency  – map of key → number of profiles it appears in
//
// TopKeys returns the most frequently occurring keys, useful for
// surfacing common configuration patterns across environments.
package stats
