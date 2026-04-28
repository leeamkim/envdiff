package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ConditionalRule defines a rule where a key must be present (and optionally
// match a value) when another key satisfies a condition.
type ConditionalRule struct {
	// WhenKey is the key whose value triggers this rule.
	WhenKey string
	// WhenValue is the value that WhenKey must have to trigger the rule.
	// An empty string means "any non-empty value".
	WhenValue string
	// ThenKey is the key that must be present when the condition is met.
	ThenKey string
	// ThenValue is the value ThenKey must have. Empty means any non-empty value.
	ThenValue string
}

// ConditionalIssue describes a violation of a ConditionalRule.
type ConditionalIssue struct {
	Rule    ConditionalRule
	EnvName string
	Reason  string
}

func (i ConditionalIssue) String() string {
	return fmt.Sprintf("[%s] %s", i.EnvName, i.Reason)
}

// CheckConditionals evaluates conditional rules against one or more env maps.
// envs maps environment name to its key/value pairs.
// It returns a list of issues for any rule violations found.
func CheckConditionals(envs map[string]map[string]string, rules []ConditionalRule) []ConditionalIssue {
	var issues []ConditionalIssue

	// Sort env names for deterministic output.
	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	for _, envName := range envNames {
		env := envs[envName]
		for _, rule := range rules {
			whenVal, whenPresent := env[rule.WhenKey]
			if !whenPresent || whenVal == "" {
				continue
			}
			// Check if the condition on WhenKey is satisfied.
			if rule.WhenValue != "" && whenVal != rule.WhenValue {
				continue
			}

			// Condition is met — now check ThenKey.
			thenVal, thenPresent := env[rule.ThenKey]
			if !thenPresent || thenVal == "" {
				issue := ConditionalIssue{
					Rule:    rule,
					EnvName: envName,
					Reason:  fmt.Sprintf("%s=%q requires %s to be set", rule.WhenKey, whenVal, rule.ThenKey),
				}
				issues = append(issues, issue)
				continue
			}
			if rule.ThenValue != "" && thenVal != rule.ThenValue {
				issue := ConditionalIssue{
					Rule:    rule,
					EnvName: envName,
					Reason:  fmt.Sprintf("%s=%q requires %s=%q, got %q", rule.WhenKey, whenVal, rule.ThenKey, rule.ThenValue, thenVal),
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

// ParseConditionalRules parses a slice of rule strings in the format:
//
//	"WHEN_KEY=WHEN_VAL:THEN_KEY" or "WHEN_KEY=WHEN_VAL:THEN_KEY=THEN_VAL"
//
// A rule like "FEATURE_X=true:FEATURE_X_URL" means: if FEATURE_X is "true",
// then FEATURE_X_URL must be present.
func ParseConditionalRules(specs []string) ([]ConditionalRule, error) {
	var rules []ConditionalRule
	for _, spec := range specs {
		parts := strings.SplitN(spec, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid conditional rule %q: expected format WHEN_KEY[=WHEN_VAL]:THEN_KEY[=THEN_VAL]", spec)
		}
		whenPart := parts[0]
		thenPart := parts[1]

		rule := ConditionalRule{}

		if idx := strings.Index(whenPart, "="); idx >= 0 {
			rule.WhenKey = whenPart[:idx]
			rule.WhenValue = whenPart[idx+1:]
		} else {
			rule.WhenKey = whenPart
		}

		if idx := strings.Index(thenPart, "="); idx >= 0 {
			rule.ThenKey = thenPart[:idx]
			rule.ThenValue = thenPart[idx+1:]
		} else {
			rule.ThenKey = thenPart
		}

		if rule.WhenKey == "" || rule.ThenKey == "" {
			return nil, fmt.Errorf("invalid conditional rule %q: WhenKey and ThenKey must not be empty", spec)
		}

		rules = append(rules, rule)
	}
	return rules, nil
}

// FormatConditionalIssues returns a human-readable summary of conditional issues.
func FormatConditionalIssues(issues []ConditionalIssue) string {
	if len(issues) == 0 {
		return "No conditional issues found."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d conditional issue(s) found:\n", len(issues)))
	for _, issue := range issues {
		sb.WriteString(fmt.Sprintf("  %s\n", issue.String()))
	}
	return strings.TrimRight(sb.String(), "\n")
}
