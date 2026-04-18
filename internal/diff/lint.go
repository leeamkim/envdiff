package diff

import "fmt"

// LintRule is a function that checks a key-value pair and returns an issue string or empty.
type LintRule func(key, value string) string

// LintIssue represents a single linting problem found in an env file.
type LintIssue struct {
	File string
	Key  string
	Msg  string
}

func (l LintIssue) String() string {
	return fmt.Sprintf("%s: [%s] %s", l.File, l.Key, l.Msg)
}

// LintResult holds all issues found across one or more files.
type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) HasIssues() bool {
	return len(r.Issues) > 0
}

// Lint runs the provided rules against each key-value pair in env and collects issues.
func Lint(file string, env map[string]string, rules []LintRule) LintResult {
	result := LintResult{}
	for key, value := range env {
		for _, rule := range rules {
			if msg := rule(key, value); msg != "" {
				result.Issues = append(result.Issues, LintIssue{
					File: file,
					Key:  key,
					Msg:  msg,
				})
			}
		}
	}
	return result
}

// LintRuleNoEmptyValues flags keys with empty values.
func LintRuleNoEmptyValues(key, value string) string {
	if value == "" {
		return "value is empty"
	}
	return ""
}

// LintRuleNoDuplicatePlaceholder flags values that look like unfilled placeholders.
func LintRuleNoDuplicatePlaceholder(key, value string) string {
	if value == "CHANGEME" || value == "TODO" || value == "PLACEHOLDER" {
		return fmt.Sprintf("value appears to be a placeholder: %q", value)
	}
	return ""
}

// LintRuleKeyUppercase flags keys that are not fully uppercase.
func LintRuleKeyUppercase(key, value string) string {
	for _, c := range key {
		if c >= 'a' && c <= 'z' {
			return "key contains lowercase letters"
		}
	}
	return ""
}
