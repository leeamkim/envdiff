package diff

import (
	"fmt"
	"sort"
	"strings"
)

// TrimIssue represents a key whose value has leading or trailing whitespace.
type TrimIssue struct {
	Key      string
	Original string
	Trimmed  string
}

func (t TrimIssue) String() string {
	return fmt.Sprintf("%s: %q -> %q", t.Key, t.Original, t.Trimmed)
}

// CheckTrim scans env for keys whose values have leading or trailing whitespace
// and returns a list of issues describing the original and trimmed values.
func CheckTrim(env map[string]string) []TrimIssue {
	var issues []TrimIssue
	for k, v := range env {
		trimmed := strings.TrimSpace(v)
		if trimmed != v {
			issues = append(issues, TrimIssue{
				Key:      k,
				Original: v,
				Trimmed:  trimmed,
			})
		}
	}
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})
	return issues
}

// FormatTrimIssues returns a human-readable summary of trim issues.
func FormatTrimIssues(issues []TrimIssue) string {
	if len(issues) == 0 {
		return "no trim issues found"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d trim issue(s) found:\n", len(issues))
	for _, issue := range issues {
		fmt.Fprintf(&sb, "  %s\n", issue.String())
	}
	return strings.TrimRight(sb.String(), "\n")
}
