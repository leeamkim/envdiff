package diff

import (
	"fmt"
	"sort"
	"strings"
)

// WhitelistIssue represents a key that violates an allowlist rule.
type WhitelistIssue struct {
	Key    string
	Value  string
	Reason string
}

func (w WhitelistIssue) String() string {
	return fmt.Sprintf("%s: %s", w.Key, w.Reason)
}

// WhitelistRule is a function that checks a key/value pair and returns an issue if violated.
type WhitelistRule func(key, value string) *WhitelistIssue

// WhitelistRuleAllowedKeys returns a rule that flags any key not in the allowed set.
func WhitelistRuleAllowedKeys(allowed []string) WhitelistRule {
	set := make(map[string]struct{}, len(allowed))
	for _, k := range allowed {
		set[strings.ToUpper(k)] = struct{}{}
	}
	return func(key, value string) *WhitelistIssue {
		if _, ok := set[strings.ToUpper(key)]; !ok {
			return &WhitelistIssue{
				Key:    key,
				Value:  value,
				Reason: "key not in allowlist",
			}
		}
		return nil
	}
}

// WhitelistRuleAllowedValues returns a rule that flags a key whose value is not in the allowed set.
func WhitelistRuleAllowedValues(key string, allowed []string) WhitelistRule {
	set := make(map[string]struct{}, len(allowed))
	for _, v := range allowed {
		set[v] = struct{}{}
	}
	upperKey := strings.ToUpper(key)
	return func(k, value string) *WhitelistIssue {
		if strings.ToUpper(k) != upperKey {
			return nil
		}
		if _, ok := set[value]; !ok {
			return &WhitelistIssue{
				Key:    k,
				Value:  value,
				Reason: fmt.Sprintf("value %q not in allowed values", value),
			}
		}
		return nil
	}
}

// CheckWhitelist applies the given rules to every key/value in env and returns all issues.
func CheckWhitelist(env map[string]string, rules []WhitelistRule) []WhitelistIssue {
	var issues []WhitelistIssue
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, rule := range rules {
			if issue := rule(k, env[k]); issue != nil {
				issues = append(issues, *issue)
			}
		}
	}
	return issues
}

// FormatWhitelistIssues returns a human-readable summary of whitelist issues.
func FormatWhitelistIssues(issues []WhitelistIssue) string {
	if len(issues) == 0 {
		return "no whitelist violations found"
	}
	var sb strings.Builder
	for _, i := range issues {
		sb.WriteString(fmt.Sprintf("  [VIOLATION] %s\n", i.String()))
	}
	return sb.String()
}
