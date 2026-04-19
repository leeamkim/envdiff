package diff

import "fmt"

// RequiredIssue represents a missing required key in an env map.
type RequiredIssue struct {
	Key string
}

func (r RequiredIssue) String() string {
	return fmt.Sprintf("missing required key: %s", r.Key)
}

// CheckRequired validates that all keys in the required list are present
// and non-empty in the provided env map.
func CheckRequired(env map[string]string, required []string) []RequiredIssue {
	var issues []RequiredIssue
	for _, key := range required {
		val, ok := env[key]
		if !ok || val == "" {
			issues = append(issues, RequiredIssue{Key: key})
		}
	}
	return issues
}

// FormatRequiredIssues returns a human-readable summary of required key issues.
func FormatRequiredIssues(issues []RequiredIssue) string {
	if len(issues) == 0 {
		return "all required keys present"
	}
	out := fmt.Sprintf("%d required key(s) missing or empty:\n", len(issues))
	for _, iss := range issues {
		out += fmt.Sprintf("  - %s\n", iss.Key)
	}
	return out
}
