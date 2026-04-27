package diff

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// EncodingIssue represents a detected encoding or character problem in an env key or value.
type EncodingIssue struct {
	Key     string
	Value   string
	Message string
}

func (e EncodingIssue) String() string {
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

// CheckEncoding scans an env map for non-ASCII characters, control characters,
// or values that appear to contain raw binary data.
func CheckEncoding(env map[string]string) []EncodingIssue {
	var issues []EncodingIssue

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]

		if hasNonASCII(k) {
			issues = append(issues, EncodingIssue{
				Key:     k,
				Value:   v,
				Message: "key contains non-ASCII characters",
			})
			continue
		}

		if hasControlChars(v) {
			issues = append(issues, EncodingIssue{
				Key:     k,
				Value:   v,
				Message: "value contains control characters",
			})
			continue
		}

		if hasNonASCII(v) {
			issues = append(issues, EncodingIssue{
				Key:     k,
				Value:   v,
				Message: "value contains non-ASCII characters",
			})
		}
	}

	return issues
}

// FormatEncodingIssues returns a human-readable summary of encoding issues.
func FormatEncodingIssues(issues []EncodingIssue) string {
	if len(issues) == 0 {
		return "no encoding issues found"
	}
	var sb strings.Builder
	for _, iss := range issues {
		sb.WriteString(fmt.Sprintf("  [encoding] %s\n", iss.String()))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func hasNonASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return true
		}
	}
	return false
}

func hasControlChars(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) && r != '\t' {
			return true
		}
	}
	return false
}
