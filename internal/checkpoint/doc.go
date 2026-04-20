// Package checkpoint provides named save-points (checkpoints) for environment
// variable profiles. A checkpoint captures the full variable map of a profile
// at a point in time, along with an optional human-readable note.
//
// Checkpoints are stored as individual JSON files in a configurable directory
// and can be listed, loaded, or deleted by name. They complement the history
// and snapshot features by giving users explicit, named restore points rather
// than automatic rolling history.
package checkpoint
