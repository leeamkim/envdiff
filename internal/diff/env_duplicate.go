package diff

import (
	"fmt"
	"sort"
	"strings"
)

// DuplicateIssue represents a key that appears with conflicting values across envs.
type DuplicateIssue struct {
	Key    string
	EnvA   string
	EnvB   string
	ValueA string
	ValueB string
}

func (d DuplicateIssue) String() string {
	return fmt.Sprintf("key %q differs: %s=%q vs %s=%q", d.Key, d.EnvA, d.ValueA, d.EnvB, d.ValueB)
}

// FindDuplicateConflicts checks for keys present in all provided envs that have
// differing values. envs is a map of env-name -> key/value pairs.
func FindDuplicateConflicts(envs map[string]map[string]string) []DuplicateIssue {
	if len(envs) < 2 {
		return nil
	}

	// Collect sorted env names for deterministic output.
	names := make([]string, 0, len(envs))
	for n := range envs {
		names = append(names, n)
	}
	sort.Strings(names)

	// Collect all keys.
	keySet := map[string]struct{}{}
	for _, m := range envs {
		for k := range m {
			keySet[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var issues []DuplicateIssue
	for _, key := range keys {
		for i := 0; i < len(names)-1; i++ {
			for j := i + 1; j < len(names); j++ {
				va, okA := envs[names[i]][key]
				vb, okB := envs[names[j]][key]
				if !okA || !okB {
					continue
				}
				if va != vb {
					issues = append(issues, DuplicateIssue{
						Key:    key,
						EnvA:   names[i],
						EnvB:   names[j],
						ValueA: va,
						ValueB: vb,
					})
				}
			}
		}
	}
	return issues
}

// FormatDuplicateIssues returns a human-readable report of duplicate conflicts.
func FormatDuplicateIssues(issues []DuplicateIssue) string {
	if len(issues) == 0 {
		return "no duplicate conflicts found"
	}
	var sb strings.Builder
	for _, iss := range issues {
		sb.WriteString(iss.String())
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}
