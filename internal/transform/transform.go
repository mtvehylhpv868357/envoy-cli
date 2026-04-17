package transform

import (
	"strings"
)

// Op represents a transformation operation.
type Op string

const (
	OpUppercase    Op = "uppercase"
	OpLowercase    Op = "lowercase"
	OpTrimSpace    Op = "trim"
	OpBase64Encode Op = "base64encode"
	OpBase64Decode Op = "base64decode"
)

// Options configures the transformation.
type Options struct {
	Keys []string // if empty, apply to all keys
	Op   Op
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Op: OpTrimSpace}
}

// Map applies the transformation to a map of env vars.
func Map(vars map[string]string, opts Options) (map[string]string, error) {
	target := buildTargetSet(opts.Keys)
	result := make(map[string]string, len(vars))
	for k, v := range vars {
		if len(target) == 0 || target[k] {
			transformed, err := applyOp(v, opts.Op)
			if err != nil {
				return nil, err
			}
			result[k] = transformed
		} else {
			result[k] = v
		}
	}
	return result, nil
}

func applyOp(val string, op Op) (string, error) {
	switch op {
	case OpUppercase:
		return strings.ToUpper(val), nil
	case OpLowercase:
		return strings.ToLower(val), nil
	case OpTrimSpace:
		return strings.TrimSpace(val), nil
	case OpBase64Encode:
		return encodeB64(val), nil
	case OpBase64Decode:
		return decodeB64(val)
	}
	return val, nil
}

func buildTargetSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
