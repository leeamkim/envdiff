package diff

import (
	"strings"
)

// NormalizeOption controls how values are normalized before comparison.
type NormalizeOption func(string) string

// NormalizeTrimSpace trims leading and trailing whitespace from values.
func NormalizeTrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// NormalizeLowercase converts values to lowercase.
func NormalizeLowercase(s string) string {
	return strings.ToLower(s)
}

// NormalizeMap applies a set of NormalizeOptions to every value in an env map,
// returning a new map with normalized values.
func NormalizeMap(env map[string]string, opts ...NormalizeOption) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		for _, opt := range opts {
			v = opt(v)
		}
		result[k] = v
	}
	return result
}

// NormalizeCompare runs Compare after normalizing both env maps with the
// provided options. Original values are not mutated.
func NormalizeCompare(a, b map[string]string, opts ...NormalizeOption) Result {
	return Compare(NormalizeMap(a, opts...), NormalizeMap(b, opts...))
}
