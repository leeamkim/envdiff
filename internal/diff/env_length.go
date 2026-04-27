package diff

import (
	"fmt"
	"sort"
	"strings"
)

// LengthIssue represents a key whose value length violates a defined rule.
type LengthIssue struct {
	Key      string
	Actual   int
	Min      int
	Max      int
	Violation string
}

func (i LengthIssue) String() string {
	return fmt.Sprintf("%s: length %d violates rule (%s, min=%d max=%d)", i.Key, i.Actual, i.Violation, i.Min, i.Max)
}

// LengthRule defines min/max length constraints for keys matching a prefix pattern.
type LengthRule struct {
	Pattern string
	Min     int
	Max     int // 0 means no upper bound
}

// CheckLengths validates env values against a set of LengthRules.
// Keys matching a rule's pattern (prefix match) are checked for value length.
func CheckLengths(env map[string]string, rules []LengthRule) []LengthIssue {
	var issues []LengthIssue

	for key, val := range env {
		for _, rule := range rules {
			if !strings.HasPrefix(key, rule.Pattern) && rule.Pattern != "*" {
				continue
			}
			l := len(val)
			if rule.Min > 0 && l < rule.Min {
				issues = append(issues, LengthIssue{
					Key:       key,
					Actual:    l,
					Min:       rule.Min,
					Max:       rule.Max,
					Violation: "too short",
				})
			} else if rule.Max > 0 && l > rule.Max {
				issues = append(issues, LengthIssue{
					Key:       key,
					Actual:    l,
					Min:       rule.Min,
					Max:       rule.Max,
					Violation: "too long",
				})
			}
			break
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})
	return issues
}

// FormatLengthIssues returns a human-readable summary of length issues.
func FormatLengthIssues(issues []LengthIssue) string {
	if len(issues) == 0 {
		return "no length issues found"
	}
	var sb strings.Builder
	for _, iss := range issues {
		sb.WriteString(iss.String())
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
