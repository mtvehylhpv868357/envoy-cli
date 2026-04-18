// Package reorder provides utilities for sorting and reordering environment
// variable keys within a profile.
//
// Three strategies are supported:
//
//   - alpha      – ascending alphabetical order (default)
//   - alpha-desc – descending alphabetical order
//   - custom     – caller-supplied key order; unlisted keys are appended
//                  alphabetically after the pinned set
//
// Example:
//
//	keys := reorder.Keys(vars, reorder.Options{Strategy: reorder.StrategyAlpha})
//	for _, k := range keys {
//		fmt.Printf("%s=%s\n", k, vars[k])
//	}
package reorder
