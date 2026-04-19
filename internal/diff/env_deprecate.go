package diff

import "fmt"

// DeprecationRule marks keys as deprecated.
type DeprecationRule struct {
	Key     string
	Reason  string
	Replace string
}

// DeprecationIssue represents a deprecated key found in an env map.
type DeprecationIssue struct {
	Key     string
	Reason  string
	Replace string
}

func (d DeprecationIssue) String() string {
	if d.Replace != "" {
		return fmt.Sprintf("DEPRECATED %s: %s (use %s instead)", d.Key, d.Reason, d.Replace)
	}
	return fmt.Sprintf("DEPRECATED %s: %s", d.Key, d.Reason)
}

// CheckDeprecations scans env for keys matching any deprecation rule.
func CheckDeprecations(env map[string]string, rules []DeprecationRule) []DeprecationIssue {
	var issues []DeprecationIssue
	for _, rule := range rules {
		if _, ok := env[rule.Key]; ok {
			issues = append(issues, DeprecationIssue{
				Key:     rule.Key,
				Reason:  rule.Reason,
				Replace: rule.Replace,
			})
		}
	}
	return issues
}

// FormatDeprecationIssues returns a human-readable summary.
func FormatDeprecationIssues(issues []DeprecationIssue) string {
	if len(issues) == 0 {
		return "no deprecated keys found"
	}
	out := ""
	for _, i := range issues {
		out += i.String() + "\n"
	}
	return out
}
