package diff

import "fmt"

// ClampRule defines a rule that checks whether a key's value length falls within [Min, Max].
type ClampRule struct {
	Key string
	Min int
	Max int
}

// ClampIssue represents a value-length violation for a specific key.
type ClampIssue struct {
	Key    string
	Value  string
	Min    int
	Max    int
	Length int
}

func (c ClampIssue) String() string {
	return fmt.Sprintf("key %q value length %d out of range [%d, %d]", c.Key, c.Length, c.Min, c.Max)
}

// CheckClamps validates that each key in env satisfies its corresponding ClampRule.
// Keys not covered by any rule are skipped.
func CheckClamps(env map[string]string, rules []ClampRule) []ClampIssue {
	ruleMap := make(map[string]ClampRule, len(rules))
	for _, r := range rules {
		ruleMap[r.Key] = r
	}

	var issues []ClampIssue
	for _, r := range rules {
		val, ok := env[r.Key]
		if !ok {
			continue
		}
		l := len(val)
		if l < r.Min || l > r.Max {
			issues = append(issues, ClampIssue{
				Key:    r.Key,
				Value:  val,
				Min:    r.Min,
				Max:    r.Max,
				Length: l,
			})
		}
	}
	return issues
}

// FormatClampIssues returns a human-readable summary of clamp violations.
func FormatClampIssues(issues []ClampIssue) string {
	if len(issues) == 0 {
		return "no clamp issues found"
	}
	out := fmt.Sprintf("%d clamp issue(s):\n", len(issues))
	for _, i := range issues {
		out += "  " + i.String() + "\n"
	}
	return out
}
