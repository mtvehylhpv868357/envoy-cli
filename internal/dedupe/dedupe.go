// Package dedupe removes duplicate environment variable entries,
// keeping the last occurrence of each key by default.
package dedupe

// Options configures deduplication behaviour.
type Options struct {
	// KeepFirst retains the first occurrence instead of the last.
	KeepFirst bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{KeepFirst: false}
}

// Map deduplicates a slice of key=value pairs and returns a clean map.
// Input order is preserved when building the result; the chosen occurrence
// (first or last) wins for each key.
func Map(pairs []string, opts Options) map[string]string {
	seen := make(map[string]string)
	order := make([]string, 0, len(pairs))

	for _, pair := range pairs {
		key, val := splitPair(pair)
		if key == "" {
			continue
		}
		if _, exists := seen[key]; !exists {
			order = append(order, key)
		}
		if opts.KeepFirst {
			if _, exists := seen[key]; !exists {
				seen[key] = val
			}
		} else {
			seen[key] = val
		}
	}

	result := make(map[string]string, len(order))
	for _, k := range order {
		result[k] = seen[k]
	}
	return result
}

// Duplicates returns only the keys that appear more than once in pairs.
func Duplicates(pairs []string) []string {
	counts := make(map[string]int)
	for _, pair := range pairs {
		key, _ := splitPair(pair)
		if key != "" {
			counts[key]++
		}
	}
	var dups []string
	seen := make(map[string]bool)
	for _, pair := range pairs {
		key, _ := splitPair(pair)
		if counts[key] > 1 && !seen[key] {
			dups = append(dups, key)
			seen[key] = true
		}
	}
	return dups
}

func splitPair(pair string) (string, string) {
	for i := 0; i < len(pair); i++ {
		if pair[i] == '=' {
			return pair[:i], pair[i+1:]
		}
	}
	return pair, ""
}
