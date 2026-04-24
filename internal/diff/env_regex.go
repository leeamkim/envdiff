package diff

import (
	"fmt"
	"regexp"
	"sort"
)

// RegexIssue represents a key whose value does not match the expected pattern.
type RegexIssue struct {
	Key     string
	Value   string
	Pattern string
}

func (r RegexIssue) String() string {
	return fmt.Sprintf("key %q value %q does not match pattern %q", r.Key, r.Value, r.Pattern)
}

// RegexRule maps a key (or glob-style prefix ending in *) to a compiled regex.
type RegexRule struct {
	Key     string
	Pattern *regexp.Regexp
}

// CheckRegex validates that env values match their declared patterns.
// rules is a slice of RegexRule; env is the key→value map to validate.
func CheckRegex(env map[string]string, rules []RegexRule) []RegexIssue {
	var issues []RegexIssue
	for _, rule := range rules {
		for key, val := range env {
			if !keyMatchesRule(key, rule.Key) {
				continue
			}
			if !rule.Pattern.MatchString(val) {
				issues = append(issues, RegexIssue{
					Key:     key,
					Value:   val,
					Pattern: rule.Pattern.String(),
				})
			}
		}
	}
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})
	return issues
}

// keyMatchesRule returns true when key equals ruleKey or ruleKey ends with
// '*' and key has the matching prefix.
func keyMatchesRule(key, ruleKey string) bool {
	if len(ruleKey) > 0 && ruleKey[len(ruleKey)-1] == '*' {
		prefix := ruleKey[:len(ruleKey)-1]
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	}
	return key == ruleKey
}

// FormatRegexIssues returns a human-readable summary of regex issues.
func FormatRegexIssues(issues []RegexIssue) string {
	if len(issues) == 0 {
		return "no regex issues found"
	}
	out := fmt.Sprintf("%d regex issue(s):\n", len(issues))
	for _, iss := range issues {
		out += "  " + iss.String() + "\n"
	}
	return out
}
