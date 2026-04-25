package diff

import (
	"fmt"
	"sort"
	"strings"
)

// AccessLevel represents the sensitivity tier of an env key.
type AccessLevel int

const (
	AccessPublic AccessLevel = iota
	AccessInternal
	AccessSecret
)

func (a AccessLevel) String() string {
	switch a {
	case AccessPublic:
		return "public"
	case AccessInternal:
		return "internal"
	case AccessSecret:
		return "secret"
	default:
		return "unknown"
	}
}

// AccessIssue describes a key whose access level violates policy.
type AccessIssue struct {
	Key      string
	Actual   AccessLevel
	Expected AccessLevel
}

func (i AccessIssue) String() string {
	return fmt.Sprintf("key %q: expected %s, got %s", i.Key, i.Expected, i.Actual)
}

// AccessRule maps a key prefix to the required AccessLevel.
type AccessRule struct {
	Prefix string
	Level  AccessLevel
}

// CheckAccess verifies that each key in env satisfies the provided rules.
// Rules are matched by key prefix (case-insensitive). The first matching rule wins.
func CheckAccess(env map[string]string, rules []AccessRule, actual map[string]AccessLevel) []AccessIssue {
	var issues []AccessIssue
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		expected, matched := matchAccessRule(k, rules)
		if !matched {
			continue
		}
		got, ok := actual[k]
		if !ok {
			got = AccessPublic
		}
		if got != expected {
			issues = append(issues, AccessIssue{Key: k, Actual: got, Expected: expected})
		}
	}
	return issues
}

func matchAccessRule(key string, rules []AccessRule) (AccessLevel, bool) {
	upper := strings.ToUpper(key)
	for _, r := range rules {
		if strings.HasPrefix(upper, strings.ToUpper(r.Prefix)) {
			return r.Level, true
		}
	}
	return AccessPublic, false
}

// FormatAccessIssues returns a human-readable summary of access issues.
func FormatAccessIssues(issues []AccessIssue) string {
	if len(issues) == 0 {
		return "no access issues found"
	}
	var sb strings.Builder
	for _, i := range issues {
		sb.WriteString("  " + i.String() + "\n")
	}
	return sb.String()
}
