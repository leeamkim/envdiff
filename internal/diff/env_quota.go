package diff

import (
	"fmt"
	"sort"
)

// QuotaRule defines a maximum number of keys allowed matching a prefix pattern.
type QuotaRule struct {
	Prefix string
	MaxKeys int
}

// QuotaIssue represents a violation of a quota rule.
type QuotaIssue struct {
	Prefix  string
	Allowed int
	Actual  int
	Keys    []string
}

func (q QuotaIssue) String() string {
	return fmt.Sprintf("prefix %q: allowed %d keys, found %d (%v)", q.Prefix, q.Allowed, q.Actual, q.Keys)
}

// CheckQuotas checks whether the number of keys matching each rule's prefix
// exceeds the allowed maximum.
func CheckQuotas(env map[string]string, rules []QuotaRule) []QuotaIssue {
	var issues []QuotaIssue

	for _, rule := range rules {
		var matched []string
		for k := range env {
			if len(k) >= len(rule.Prefix) && k[:len(rule.Prefix)] == rule.Prefix {
				matched = append(matched, k)
			}
		}
		sort.Strings(matched)
		if len(matched) > rule.MaxKeys {
			issues = append(issues, QuotaIssue{
				Prefix:  rule.Prefix,
				Allowed: rule.MaxKeys,
				Actual:  len(matched),
				Keys:    matched,
			})
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Prefix < issues[j].Prefix
	})
	return issues
}

// FormatQuotaIssues returns a human-readable summary of quota violations.
func FormatQuotaIssues(issues []QuotaIssue) string {
	if len(issues) == 0 {
		return "no quota violations found"
	}
	out := fmt.Sprintf("%d quota violation(s):\n", len(issues))
	for _, iss := range issues {
		out += "  " + iss.String() + "\n"
	}
	return out
}
