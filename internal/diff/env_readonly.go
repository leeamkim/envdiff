package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ReadonlyIssue represents a key that was modified despite being marked readonly.
type ReadonlyIssue struct {
	Key      string
	Env      string
	OldValue string
	NewValue string
}

func (i ReadonlyIssue) String() string {
	return fmt.Sprintf("[%s] %s is readonly (was %q, now %q)", i.Env, i.Key, i.OldValue, i.NewValue)
}

// CheckReadonly verifies that keys marked as readonly have not changed between
// a reference env and one or more target envs.
func CheckReadonly(reference map[string]string, targets map[string]map[string]string, readonlyKeys []string) []ReadonlyIssue {
	if len(readonlyKeys) == 0 || len(targets) == 0 {
		return nil
	}

	keySet := make(map[string]struct{}, len(readonlyKeys))
	for _, k := range readonlyKeys {
		keySet[strings.TrimSpace(k)] = struct{}{}
	}

	var issues []ReadonlyIssue

	envNames := make([]string, 0, len(targets))
	for name := range targets {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	for _, envName := range envNames {
		envVars := targets[envName]
		for key := range keySet {
			refVal, refOk := reference[key]
			if !refOk {
				continue
			}
			targetVal, targetOk := envVars[key]
			if !targetOk {
				continue
			}
			if refVal != targetVal {
				issues = append(issues, ReadonlyIssue{
					Key:      key,
					Env:      envName,
					OldValue: refVal,
					NewValue: targetVal,
				})
			}
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Env != issues[j].Env {
			return issues[i].Env < issues[j].Env
		}
		return issues[i].Key < issues[j].Key
	})

	return issues
}

// FormatReadonlyIssues formats a slice of ReadonlyIssue for display.
func FormatReadonlyIssues(issues []ReadonlyIssue) string {
	if len(issues) == 0 {
		return "no readonly violations found"
	}
	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString(issue.String())
		sb.WriteString("\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
