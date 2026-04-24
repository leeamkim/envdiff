package diff

import (
	"fmt"
	"sort"
	"strings"
)

// PrefixIssue represents a key that violates a required prefix rule.
type PrefixIssue struct {
	Key      string
	Expected string
}

func (p PrefixIssue) String() string {
	return fmt.Sprintf("key %q does not have required prefix %q", p.Key, p.Expected)
}

// CheckPrefixes validates that all keys in env start with one of the allowed prefixes.
// If prefixes is empty, no issues are returned.
func CheckPrefixes(env map[string]string, prefixes []string) []PrefixIssue {
	if len(prefixes) == 0 {
		return nil
	}

	var issues []PrefixIssue
	for key := range env {
		matched := false
		for _, p := range prefixes {
			if strings.HasPrefix(key, p) {
				matched = true
				break
			}
		}
		if !matched {
			issues = append(issues, PrefixIssue{Key: key, Expected: strings.Join(prefixes, "|")})
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})
	return issues
}

// FormatPrefixIssues returns a human-readable summary of prefix issues.
func FormatPrefixIssues(issues []PrefixIssue) string {
	if len(issues) == 0 {
		return "no prefix issues found"
	}
	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString("  " + issue.String() + "\n")
	}
	return sb.String()
}
