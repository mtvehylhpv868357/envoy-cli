package reorder

import (
	"fmt"
	"sort"

	"github.com/envoy-cli/internal/profile"
)

// Strategy controls how keys are reordered.
type Strategy string

const (
	StrategyAlpha    Strategy = "alpha"
	StrategyAlphaDesc Strategy = "alpha-desc"
	StrategyCustom   Strategy = "custom"
)

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Strategy: StrategyAlpha,
	}
}

// Options configures reorder behaviour.
type Options struct {
	Strategy Strategy
	// Order is used when Strategy == StrategyCustom; keys not listed appear last.
	Order []string
}

// Profile reorders the keys of the named profile and saves it.
func Profile(st *profile.Store, name string, opts Options) (map[string]string, error) {
	if name == "" {
		return nil, fmt.Errorf("profile name must not be empty")
	}
	vars, err := st.Get(name)
	if err != nil {
		return nil, fmt.Errorf("profile %q not found: %w", name, err)
	}

	ordered := Map(vars, opts)

	if err := st.Add(name, ordered); err != nil {
		return nil, fmt.Errorf("saving profile: %w", err)
	}
	return ordered, nil
}

// Map returns a new map with keys ordered according to opts.
// Because Go maps are unordered the result is returned alongside a
// canonical key slice so callers can iterate in order.
func Map(vars map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		out[k] = v
	}
	return out
}

// Keys returns the keys of vars in the order described by opts.
func Keys(vars map[string]string, opts Options) []string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}

	switch opts.Strategy {
	case StrategyAlphaDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	case StrategyCustom:
		keys = customOrder(keys, opts.Order)
	default: // StrategyAlpha
		sort.Strings(keys)
	}
	return keys
}

func customOrder(keys []string, order []string) []string {
	rank := make(map[string]int, len(order))
	for i, k := range order {
		rank[k] = i
	}

	pinned := []string{}
	rest := []string{}
	for _, k := range keys {
		if _, ok := rank[k]; ok {
			pinned = append(pinned, k)
		} else {
			rest = append(rest, k)
		}
	}
	sort.Slice(pinned, func(i, j int) bool {
		return rank[pinned[i]] < rank[pinned[j]]
	})
	sort.Strings(rest)
	return append(pinned, rest...)
}
