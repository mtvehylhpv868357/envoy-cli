// Package envsign provides HMAC-SHA256 signing and verification for
// envoy-cli environment variable profiles.
//
// It allows users to sign a profile's key-value pairs with a passphrase,
// producing a SignedEnvelope that can later be verified to detect tampering.
//
// Usage:
//
//	env, err := envsign.Sign("production", vars, "my-secret")
//	if err != nil { ... }
//
//	if err := envsign.Verify(env, "my-secret"); err != nil {
//		// signature mismatch or tampering detected
//	}
package envsign
