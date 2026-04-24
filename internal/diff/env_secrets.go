package diff

import (
	"fmt"
	"sort"
	"strings"
)

// SecretIssue represents a key that appears to contain a weak or default secret.
type SecretIssue struct {
	Key    string
	Value  string
	Reason string
}

func (s SecretIssue) String() string {
	return fmt.Sprintf("[SECRET] %s: %s", s.Key, s.Reason)
}

// SecretRule is a function that evaluates a key/value pair for secret hygiene.
type SecretRule func(key, value string) (string, bool)

// SecretRuleWeakValue flags values that are common weak secrets.
func SecretRuleWeakValue(key, value string) (string, bool) {
	weak := []string{"secret", "password", "changeme", "12345", "admin", "letmein", "test"}
	lower := strings.ToLower(value)
	for _, w := range weak {
		if lower == w {
			return fmt.Sprintf("value matches known weak secret %q", w), true
		}
	}
	return "", false
}

// SecretRuleKeyWithoutValue flags keys that look like secrets but have empty values.
func SecretRuleKeyWithoutValue(key, value string) (string, bool) {
	sensitivePatterns := []string{"secret", "password", "token", "key", "api", "auth"}
	lowerKey := strings.ToLower(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(lowerKey, p) && strings.TrimSpace(value) == "" {
			return "sensitive key has empty value", true
		}
	}
	return "", false
}

// CheckSecrets runs secret hygiene rules against an env map.
func CheckSecrets(env map[string]string, rules []SecretRule) []SecretIssue {
	var issues []SecretIssue
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		for _, rule := range rules {
			if reason, ok := rule(k, v); ok {
				issues = append(issues, SecretIssue{Key: k, Value: v, Reason: reason})
				break
			}
		}
	}
	return issues
}

// FormatSecretIssues returns a human-readable summary of secret issues.
func FormatSecretIssues(issues []SecretIssue) string {
	if len(issues) == 0 {
		return "no secret issues found"
	}
	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString(issue.String())
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
