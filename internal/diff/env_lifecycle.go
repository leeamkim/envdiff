package diff

import (
	"fmt"
	"sort"
	"strings"
)

// LifecycleStage represents the maturity stage of an env key.
type LifecycleStage string

const (
	StageActive     LifecycleStage = "active"
	StageDeprecated LifecycleStage = "deprecated"
	StageExperimental LifecycleStage = "experimental"
	StageRetired    LifecycleStage = "retired"
)

// LifecycleRule maps a key pattern to a lifecycle stage.
type LifecycleRule struct {
	Pattern string
	Stage   LifecycleStage
}

// LifecycleIssue describes a key whose stage violates policy.
type LifecycleIssue struct {
	Key   string
	Stage LifecycleStage
	Note  string
}

func (i LifecycleIssue) String() string {
	return fmt.Sprintf("%s [%s]: %s", i.Key, i.Stage, i.Note)
}

// CheckLifecycle inspects env keys against lifecycle rules and flags
// keys that are deprecated or retired.
func CheckLifecycle(env map[string]string, rules []LifecycleRule) []LifecycleIssue {
	var issues []LifecycleIssue

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		for _, rule := range rules {
			if matchesLifecyclePattern(key, rule.Pattern) {
				switch rule.Stage {
				case StageDeprecated:
					issues = append(issues, LifecycleIssue{
						Key:   key,
						Stage: rule.Stage,
						Note:  "key is deprecated and should be migrated",
					})
				case StageRetired:
					issues = append(issues, LifecycleIssue{
						Key:   key,
						Stage: rule.Stage,
						Note:  "key is retired and must be removed",
					})
				}
				break
			}
		}
	}
	return issues
}

func matchesLifecyclePattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
	}
	return key == pattern
}

// FormatLifecycleIssues returns a human-readable report of lifecycle issues.
func FormatLifecycleIssues(issues []LifecycleIssue) string {
	if len(issues) == 0 {
		return "no lifecycle issues found"
	}
	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString(issue.String())
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}
