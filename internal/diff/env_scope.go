package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ScopeIssue represents a key that violates scope rules.
type ScopeIssue struct {
	Key     string
	Env     string
	Scope   string
	Message string
}

func (s ScopeIssue) String() string {
	return fmt.Sprintf("[%s] key %q violates scope %q: %s", s.Env, s.Key, s.Scope, s.Message)
}

// ScopeRule defines a named scope with an allowed key prefix set.
type ScopeRule struct {
	Scope    string
	Prefixes []string
}

// CheckScopes verifies that keys in each env conform to the declared scope rules.
// Keys not matching any allowed prefix for the given scope are flagged.
func CheckScopes(envs map[string]map[string]string, rules []ScopeRule) []ScopeIssue {
	if len(rules) == 0 || len(envs) == 0 {
		return nil
	}

	var issues []ScopeIssue

	for _, rule := range rules {
		env, ok := envs[rule.Scope]
		if !ok {
			continue
		}
		keys := make([]string, 0, len(env))
		for k := range env {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			if !matchesAnyPrefix(key, rule.Prefixes) {
				issues = append(issues, ScopeIssue{
					Key:     key,
					Env:     rule.Scope,
					Scope:   rule.Scope,
					Message: fmt.Sprintf("does not match allowed prefixes: %s", strings.Join(rule.Prefixes, ", ")),
				})
			}
		}
	}

	return issues
}

func matchesAnyPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

// FormatScopeIssues returns a human-readable summary of scope violations.
func FormatScopeIssues(issues []ScopeIssue) string {
	if len(issues) == 0 {
		return "no scope violations found"
	}
	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString(issue.String())
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}
