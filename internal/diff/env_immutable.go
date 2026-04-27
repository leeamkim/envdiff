package diff

import (
	"fmt"
	"sort"
)

// ImmutableIssue represents a key that was changed despite being declared immutable.
type ImmutableIssue struct {
	Key      string
	EnvA     string
	EnvB     string
	ValueA   string
	ValueB   string
}

func (i ImmutableIssue) String() string {
	return fmt.Sprintf("immutable key %q changed: %s=%q -> %s=%q", i.Key, i.EnvA, i.ValueA, i.EnvB, i.ValueB)
}

// CheckImmutable verifies that keys declared as immutable have the same value
// across all provided envs. envs is a map of envName -> key/value pairs.
func CheckImmutable(envs map[string]map[string]string, immutableKeys []string) []ImmutableIssue {
	if len(immutableKeys) == 0 || len(envs) < 2 {
		return nil
	}

	// Stable ordering of env names for deterministic pair comparison.
	names := make([]string, 0, len(envs))
	for n := range envs {
		names = append(names, n)
	}
	sort.Strings(names)

	var issues []ImmutableIssue

	for _, key := range immutableKeys {
		for i := 0; i < len(names)-1; i++ {
			for j := i + 1; j < len(names); j++ {
				nA, nB := names[i], names[j]
				vA, okA := envs[nA][key]
				vB, okB := envs[nB][key]
				if !okA || !okB {
					continue
				}
				if vA != vB {
					issues = append(issues, ImmutableIssue{
						Key: key, EnvA: nA, EnvB: nB, ValueA: vA, ValueB: vB,
					})
				}
			}
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Key != issues[j].Key {
			return issues[i].Key < issues[j].Key
		}
		return issues[i].EnvA < issues[j].EnvA
	})
	return issues
}

// FormatImmutableIssues returns a human-readable summary of immutable violations.
func FormatImmutableIssues(issues []ImmutableIssue) string {
	if len(issues) == 0 {
		return "no immutable violations found"
	}
	out := fmt.Sprintf("%d immutable violation(s):\n", len(issues))
	for _, iss := range issues {
		out += "  " + iss.String() + "\n"
	}
	return out
}
