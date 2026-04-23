package diff

import "fmt"

// DefaultIssue represents a key whose value matches the default placeholder.
type DefaultIssue struct {
	Key          string
	Value        string
	DefaultValue string
}

func (d DefaultIssue) String() string {
	return fmt.Sprintf("key %q has default value %q", d.Key, d.DefaultValue)
}

// DefaultRule is a function that checks whether a value is considered a default.
type DefaultRule func(key, value string) (defaultVal string, matched bool)

// DefaultRuleCommonPlaceholders flags values like "changeme", "todo", "default", "example".
func DefaultRuleCommonPlaceholders(key, value string) (string, bool) {
	common := []string{"changeme", "todo", "default", "example", "replace_me", "your_value_here"}
	lower := toLower(value)
	for _, p := range common {
		if lower == p {
			return p, true
		}
	}
	return "", false
}

// DefaultRulePinnedValues checks whether a key's value matches a pinned default map.
func DefaultRulePinnedValues(pins map[string]string) DefaultRule {
	return func(key, value string) (string, bool) {
		if def, ok := pins[key]; ok && value == def {
			return def, true
		}
		return "", false
	}
}

// CheckDefaults inspects env for keys whose values match any of the provided rules.
func CheckDefaults(env map[string]string, rules []DefaultRule) []DefaultIssue {
	var issues []DefaultIssue
	for _, k := range sortedStringKeys(env) {
		v := env[k]
		for _, rule := range rules {
			if def, matched := rule(k, v); matched {
				issues = append(issues, DefaultIssue{Key: k, Value: v, DefaultValue: def})
				break
			}
		}
	}
	return issues
}

// FormatDefaultIssues returns a human-readable list of default issues.
func FormatDefaultIssues(issues []DefaultIssue) string {
	if len(issues) == 0 {
		return "no default value issues found"
	}
	out := ""
	for _, iss := range issues {
		out += fmt.Sprintf("  [DEFAULT] %s\n", iss.String())
	}
	return out
}

func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

func sortedStringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
