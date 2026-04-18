package diff

import "fmt"

// SchemaIssue represents a violation found during schema validation.
type SchemaIssue struct {
	Key     string
	Message string
}

func (s SchemaIssue) String() string {
	return fmt.Sprintf("[%s] %s", s.Key, s.Message)
}

// SchemaRule is a function that validates a key-value pair against a schema.
type SchemaRule func(key, value string) *SchemaIssue

// ValidateSchema runs all provided rules against the given env map and returns issues.
func ValidateSchema(env map[string]string, rules []SchemaRule) []SchemaIssue {
	var issues []SchemaIssue
	for key, value := range env {
		for _, rule := range rules {
			if issue := rule(key, value); issue != nil {
				issues = append(issues, *issue)
			}
		}
	}
	return issues
}

// SchemaRuleRequiredKeys returns a rule that flags keys missing from the env map.
func SchemaRuleRequiredKeys(required []string) func(env map[string]string) []SchemaIssue {
	return func(env map[string]string) []SchemaIssue {
		var issues []SchemaIssue
		for _, key := range required {
			if _, ok := env[key]; !ok {
				issues = append(issues, SchemaIssue{Key: key, Message: "required key is missing"})
			}
		}
		return issues
	}
}

// SchemaRuleNoEmptyValues flags keys with empty values.
func SchemaRuleNoEmptyValues(key, value string) *SchemaIssue {
	if value == "" {
		return &SchemaIssue{Key: key, Message: "value must not be empty"}
	}
	return nil
}
