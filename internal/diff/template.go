package diff

import (
	"fmt"
	"sort"
	"strings"
)

// TemplateIssue represents a missing or extra key relative to a template.
type TemplateIssue struct {
	Key    string
	Reason string
}

func (t TemplateIssue) String() string {
	return fmt.Sprintf("%s: %s", t.Key, t.Reason)
}

// CompareToTemplate checks that env contains all keys in template and no extra keys
// (if strict is true).
func CompareToTemplate(template, env map[string]string, strict bool) []TemplateIssue {
	var issues []TemplateIssue

	for key := range template {
		if _, ok := env[key]; !ok {
			issues = append(issues, TemplateIssue{Key: key, Reason: "missing from env"})
		}
	}

	if strict {
		for key := range env {
			if _, ok := template[key]; !ok {
				issues = append(issues, TemplateIssue{Key: key, Reason: "extra key not in template"})
			}
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})
	return issues
}

// FormatTemplateIssues returns a human-readable summary.
func FormatTemplateIssues(issues []TemplateIssue) string {
	if len(issues) == 0 {
		return "no template issues found"
	}
	var sb strings.Builder
	for _, iss := range issues {
		fmt.Fprintf(&sb, "  [template] %s\n", iss)
	}
	return strings.TrimRight(sb.String(), "\n")
}
