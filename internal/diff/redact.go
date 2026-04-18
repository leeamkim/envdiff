package diff

import "strings"

// RedactOptions controls which keys are redacted in output.
type RedactOptions struct {
	Patterns []string // substrings to match against key names (case-insensitive)
}

// DefaultRedactPatterns are common sensitive key substrings.
var DefaultRedactPatterns = []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL"}

// Redact replaces values of sensitive keys with "***" in the given env map.
func Redact(env map[string]string, opts RedactOptions) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, opts.Patterns) {
			result[k] = "***"
		} else {
			result[k] = v
		}
	}
	return result
}

// RedactResult replaces sensitive values inside a Result's diff entries.
func RedactResult(r Result, opts RedactOptions) Result {
	out := Result{
		MissingInA:  make(map[string]string),
		MissingInB:  make(map[string]string),
		Mismatched:  make(map[string][2]string),
	}
	for k, v := range r.MissingInA {
		if isSensitive(k, opts.Patterns) {
			out.MissingInA[k] = "***"
		} else {
			out.MissingInA[k] = v
		}
	}
	for k, v := range r.MissingInB {
		if isSensitive(k, opts.Patterns) {
			out.MissingInB[k] = "***"
		} else {
			out.MissingInB[k] = v
		}
	}
	for k, pair := range r.Mismatched {
		if isSensitive(k, opts.Patterns) {
			out.Mismatched[k] = [2]string{"***", "***"}
		} else {
			out.Mismatched[k] = pair
		}
	}
	return out
}

func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
