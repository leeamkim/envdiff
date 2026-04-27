package diff

import (
	"fmt"
	"sort"
	"strings"
)

// CrossRefIssue represents a cross-reference problem between two env maps.
type CrossRefIssue struct {
	Key      string
	RefKey   string
	EnvName  string
	Reason   string
}

func (i CrossRefIssue) String() string {
	return fmt.Sprintf("%s: key %q references %q in %s — %s", i.EnvName, i.Key, i.RefKey, i.EnvName, i.Reason)
}

// CrossRefRule defines a dependency between two keys across envs.
type CrossRefRule struct {
	// SourceKey must be present and non-empty when TargetKey is set.
	SourceKey string
	TargetKey string
}

// CheckCrossRefs verifies that when a key references another key's value
// (via ${REF} syntax or explicit rules), those references are satisfiable.
func CheckCrossRefs(env map[string]string, rules []CrossRefRule, envName string) []CrossRefIssue {
	var issues []CrossRefIssue

	// Check explicit rules first.
	for _, rule := range rules {
		if val, ok := env[rule.SourceKey]; ok && val != "" {
			if refVal, exists := env[rule.TargetKey]; !exists || refVal == "" {
				issues = append(issues, CrossRefIssue{
					Key:     rule.SourceKey,
					RefKey:  rule.TargetKey,
					EnvName: envName,
					Reason:  "referenced key is missing or empty",
				})
			}
		}
	}

	// Check inline ${REF} references in values.
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		refs := extractCrossRefs(v)
		for _, ref := range refs {
			if refVal, exists := env[ref]; !exists || refVal == "" {
				issues = append(issues, CrossRefIssue{
					Key:     k,
					RefKey:  ref,
					EnvName: envName,
					Reason:  "inline reference ${" + ref + "} cannot be resolved",
				})
			}
		}
	}

	return issues
}

// extractCrossRefs finds all ${KEY} references in a value string.
func extractCrossRefs(value string) []string {
	var refs []string
	remaining := value
	for {
		start := strings.Index(remaining, "${")
		if start == -1 {
			break
		}
		end := strings.Index(remaining[start:], "}")
		if end == -1 {
			break
		}
		ref := remaining[start+2 : start+end]
		if ref != "" {
			refs = append(refs, ref)
		}
		remaining = remaining[start+end+1:]
	}
	return refs
}

// FormatCrossRefIssues returns a human-readable summary of cross-ref issues.
func FormatCrossRefIssues(issues []CrossRefIssue) string {
	if len(issues) == 0 {
		return "no cross-reference issues found"
	}
	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString(issue.String())
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
