// Package pivot transposes environment variable profiles into a
// key-centric view, making it easy to compare values for the same
// key across multiple profiles side by side.
//
// Example usage:
//
//	profiles := map[string]map[string]string{
//	    "dev":  {"HOST": "localhost"},
//	    "prod": {"HOST": "example.com"},
//	}
//	rows := pivot.Profiles(profiles, pivot.DefaultOptions())
//	for _, row := range rows {
//	    fmt.Println(row.Key, row.Values)
//	}
package pivot
