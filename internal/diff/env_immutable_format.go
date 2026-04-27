package diff

import (
	"fmt"
	"strings"
)

// SummarizeImmutable returns a one-line summary of immutable check results.
func SummarizeImmutable(issues []ImmutableIssue) string {
	if len(issues) == 0 {
		return "immutable: OK"
	}
	keys := uniqueImmutableKeys(issues)
	return fmt.Sprintf("immutable: %d violation(s) across %d key(s): %s",
		len(issues), len(keys), strings.Join(keys, ", "))
}

// FormatImmutableReport returns a detailed multi-line report grouped by key.
func FormatImmutableReport(issues []ImmutableIssue) string {
	if len(issues) == 0 {
		return "No immutable violations.\n"
	}

	grouped := make(map[string][]ImmutableIssue)
	for _, iss := range issues {
		grouped[iss.Key] = append(grouped[iss.Key], iss)
	}

	keys := uniqueImmutableKeys(issues)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Immutable Violations (%d):\n", len(issues)))
	for _, key := range keys {
		sb.WriteString(fmt.Sprintf("  [%s]\n", key))
		for _, iss := range grouped[key] {
			sb.WriteString(fmt.Sprintf("    %s=%q  vs  %s=%q\n",
				iss.EnvA, iss.ValueA, iss.EnvB, iss.ValueB))
		}
	}
	return sb.String()
}

func uniqueImmutableKeys(issues []ImmutableIssue) []string {
	seen := make(map[string]struct{})
	var keys []string
	for _, iss := range issues {
		if _, ok := seen[iss.Key]; !ok {
			seen[iss.Key] = struct{}{}
			keys = append(keys, iss.Key)
		}
	}
	return keys
}
