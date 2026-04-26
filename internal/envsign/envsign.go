// Package envsign provides HMAC-based signing and verification for
// environment variable profiles, ensuring integrity of stored profiles.
package envsign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

// ErrInvalidSignature is returned when a signature does not match.
var ErrInvalidSignature = errors.New("envsign: invalid signature")

// SignedEnvelope wraps a profile's vars with a signature.
type SignedEnvelope struct {
	Profile   string            `json:"profile"`
	Vars      map[string]string `json:"vars"`
	Signature string            `json:"signature"`
}

// Sign creates a SignedEnvelope for the given profile name and vars using
// the provided passphrase as the HMAC key.
func Sign(profile string, vars map[string]string, passphrase string) (*SignedEnvelope, error) {
	if passphrase == "" {
		return nil, errors.New("envsign: passphrase must not be empty")
	}
	sig, err := computeSignature(profile, vars, passphrase)
	if err != nil {
		return nil, fmt.Errorf("envsign: sign: %w", err)
	}
	return &SignedEnvelope{
		Profile:   profile,
		Vars:      vars,
		Signature: sig,
	}, nil
}

// Verify checks that the envelope's signature matches the given passphrase.
// Returns ErrInvalidSignature if the signature is wrong.
func Verify(env *SignedEnvelope, passphrase string) error {
	expected, err := computeSignature(env.Profile, env.Vars, passphrase)
	if err != nil {
		return fmt.Errorf("envsign: verify: %w", err)
	}
	if !hmac.Equal([]byte(expected), []byte(env.Signature)) {
		return ErrInvalidSignature
	}
	return nil
}

// computeSignature builds a deterministic canonical form of the vars and
// computes HMAC-SHA256 over profile+canonical.
func computeSignature(profile string, vars map[string]string, passphrase string) (string, error) {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ordered := make([]string, 0, len(keys))
	for _, k := range keys {
		ordered = append(ordered, k+"="+vars[k])
	}

	payload := map[string]interface{}{
		"profile": profile,
		"vars":    ordered,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, []byte(passphrase))
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil)), nil
}
