package transform

import (
	"encoding/base64"
	"fmt"
)

func encodeB64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func decodeB64(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	return string(b), nil
}
