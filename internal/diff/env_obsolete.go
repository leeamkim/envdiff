package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ObsoleteIssue represents a key that exists in the env but not in the reference.
type ObsoleteIssue struct {
	Key string
	Value string
}

func (o ObsoleteIssue) String() string {
	return fmt.Sprintf("obsolete key %q (value: %q)", o.Key, o.Value)
}

// CheckObsolete returns keys present in env but absent from the reference set.
// The reference is treated as the canonical set of allowed keys.
func CheckObsolete(env map[string]string, reference map[string]string) []ObsoleteIssue {
	var issues []ObsoleteIssue
	for k, v := range env {
		if _, ok := reference[k]; !ok {
			issues = append(issues, ObsoleteIssue{Key: k, Value: v})
		}
	}
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})
	return issues
}

// FormatObsoleteIssues returns a human-readable report of obsolete keys.
func FormatObsoleteIssues(issues []ObsoleteIssue) string {
	if len(issues) == 0 {
		return "no obsolete keys found"
	}
	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString("  " + issue.String() + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
