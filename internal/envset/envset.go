// Package envset provides set-theoretic operations on environment variable maps.
package envset

// DefaultOptions returns the default options for set operations.
func DefaultOptions() Options {
	return Options{
		CaseSensitive: true,
	}
}

// Options controls behavior of set operations.
type Options struct {
	CaseSensitive bool
}

// Union merges all maps, with later maps overwriting earlier ones.
func Union(opts Options, maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			key := normalizeKey(k, opts.CaseSensitive)
			result[key] = v
		}
	}
	return result
}

// Intersect returns keys present in ALL provided maps, using values from the first map.
func Intersect(opts Options, maps ...map[string]string) map[string]string {
	if len(maps) == 0 {
		return map[string]string{}
	}
	result := make(map[string]string)
	for k, v := range maps[0] {
		key := normalizeKey(k, opts.CaseSensitive)
		inAll := true
		for _, m := range maps[1:] {
			if _, ok := m[normalizeKey(k, opts.CaseSensitive)]; !ok {
				inAll = false
				break
			}
		}
		if inAll {
			result[key] = v
		}
	}
	return result
}

// Difference returns keys in base that are NOT present in any of the others.
func Difference(opts Options, base map[string]string, others ...map[string]string) map[string]string {
	result := make(map[string]string)
	exclude := make(map[string]struct{})
	for _, m := range others {
		for k := range m {
			exclude[normalizeKey(k, opts.CaseSensitive)] = struct{}{}
		}
	}
	for k, v := range base {
		if _, found := exclude[normalizeKey(k, opts.CaseSensitive)]; !found {
			result[normalizeKey(k, opts.CaseSensitive)] = v
		}
	}
	return result
}

// SymmetricDiff returns keys present in exactly one of the two maps.
func SymmetricDiff(opts Options, a, b map[string]string) map[string]string {
	onlyA := Difference(opts, a, b)
	onlyB := Difference(opts, b, a)
	return Union(opts, onlyA, onlyB)
}

func normalizeKey(k string, caseSensitive bool) string {
	if caseSensitive {
		return k
	}
	result := make([]byte, len(k))
	for i := 0; i < len(k); i++ {
		c := k[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		result[i] = c
	}
	return string(result)
}
