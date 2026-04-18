package diff

import "fmt"

// ValidationIssue represents a single validation problem found in an env map.
type ValidationIssue struct {
	Key     string
	Message string
}

func (v ValidationIssue) String() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// ValidationRule is a function that checks a key/value pair and returns an issue if invalid.
type ValidationRule func(key, value string) *ValidationIssue

// Validate runs all provided rules against every key in the env map and returns any issues found.
func Validate(env map[string]string, rules []ValidationRule) []ValidationIssue {
	var issues []ValidationIssue
	for _, key := range sortedKeys(env) {
		val := env[key]
		for _, rule := range rules {
			if issue := rule(key, val); issue != nil {
				issues = append(issues, *issue)
			}
		}
	}
	return issues
}

// RuleNoEmptyValues flags keys whose value is an empty string.
func RuleNoEmptyValues(key, value string) *ValidationIssue {
	if value == "" {
		return &ValidationIssue{Key: key, Message: "value is empty"}
	}
	return nil
}

// RuleNoWhitespaceValues flags keys whose value contains only whitespace.
func RuleNoWhitespaceValues(key, value string) *ValidationIssue {
	for _, ch := range value {
		if ch != ' ' && ch != '\t' {
			return nil
		}
	}
	if len(value) > 0 {
		return &ValidationIssue{Key: key, Message: "value contains only whitespace"}
	}
	return nil
}
