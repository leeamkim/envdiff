package diff

import (
	"fmt"
	"sort"
	"strings"
)

// DependencyIssue represents a key that references another key which is missing or empty.
type DependencyIssue struct {
	Key    string
	RefKey string
	Reason string
}

func (d DependencyIssue) String() string {
	return fmt.Sprintf("key %q depends on %q: %s", d.Key, d.RefKey, d.Reason)
}

// DependencyRule declares that Key depends on RefKey being present and non-empty.
type DependencyRule struct {
	Key    string
	RefKey string
}

// CheckDependencies verifies that for each declared rule, if Key is present
// in env then RefKey must also be present and non-empty.
func CheckDependencies(env map[string]string, rules []DependencyRule) []DependencyIssue {
	var issues []DependencyIssue
	for _, rule := range rules {
		val, hasKey := env[rule.Key]
		if !hasKey || val == "" {
			continue
		}
		refVal, hasRef := env[rule.RefKey]
		if !hasRef {
			issues = append(issues, DependencyIssue{
				Key:    rule.Key,
				RefKey: rule.RefKey,
				Reason: "dependency key is missing",
			})
		} else if strings.TrimSpace(refVal) == "" {
			issues = append(issues, DependencyIssue{
				Key:    rule.Key,
				RefKey: rule.RefKey,
				Reason: "dependency key is empty",
			})
		}
	}
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Key != issues[j].Key {
			return issues[i].Key < issues[j].Key
		}
		return issues[i].RefKey < issues[j].RefKey
	})
	return issues
}

// FormatDependencyIssues formats a slice of DependencyIssue into human-readable lines.
func FormatDependencyIssues(issues []DependencyIssue) string {
	if len(issues) == 0 {
		return "no dependency issues found"
	}
	var sb strings.Builder
	for _, iss := range issues {
		sb.WriteString("  " + iss.String() + "\n")
	}
	return sb.String()
}
