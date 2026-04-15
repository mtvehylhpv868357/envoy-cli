// Package compare implements side-by-side comparison of two environment
// profiles, categorising each variable as only in profile A, only in profile B,
// different between the two, or identical in both.
//
// Typical usage:
//
//	result, err := compare.Profiles(store, "dev", "prod")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, key := range result.AllKeys() {
//	    // render key according to its category
//	}
package compare
