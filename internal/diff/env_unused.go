package diff

import (
	"fmt"
	"sort"
)

// UnusedIssue represents a key present in the reference env but absent in all target envs.
type UnusedIssue struct {
	Key string
}

func (u UnusedIssue) String() string {
	return fmt.Sprintf("unused key: %s", u.Key)
}

// CheckUnused returns keys that exist in reference but are missing from every
// one of the provided target environments.
func CheckUnused(reference map[string]string, targets []map[string]string) []UnusedIssue {
	if len(targets) == 0 {
		return nil
	}

	var issues []UnusedIssue
	for key := range reference {
		foundInAny := false
		for _, t := range targets {
			if _, ok := t[key]; ok {
				foundInAny = true
				break
			}
		}
		if !foundInAny {
			issues = append(issues, UnusedIssue{Key: key})
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})
	return issues
}

// FormatUnusedIssues returns a human-readable report of unused keys.
func FormatUnusedIssues(issues []UnusedIssue) string {
	if len(issues) == 0 {
		return "no unused keys found\n"
	}
	out := fmt.Sprintf("%d unused key(s):\n", len(issues))
	for _, iss := range issues {
		out += fmt.Sprintf("  - %s\n", iss.Key)
	}
	return out
}
