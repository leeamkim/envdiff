package diff

import (
	"fmt"
	"sort"
	"strings"
)

// SuggestIssue represents a suggestion for a key that may be misnamed or misplaced.
type SuggestIssue struct {
	Key        string
	Suggestion string
	Reason     string
}

func (s SuggestIssue) String() string {
	return fmt.Sprintf("%s -> %s (%s)", s.Key, s.Suggestion, s.Reason)
}

// SuggestRenames compares keys in env against a reference set and suggests
// possible renames based on case-folded or underscore-normalized similarity.
func SuggestRenames(env map[string]string, reference map[string]string) []SuggestIssue {
	var issues []SuggestIssue

	refKeys := make([]string, 0, len(reference))
	for k := range reference {
		refKeys = append(refKeys, k)
	}

	for key := range env {
		if _, ok := reference[key]; ok {
			continue
		}
		if suggestion, reason, ok := findSuggestion(key, refKeys); ok {
			issues = append(issues, SuggestIssue{
				Key:        key,
				Suggestion: suggestion,
				Reason:     reason,
			})
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})
	return issues
}

func findSuggestion(key string, refKeys []string) (string, string, bool) {
	normKey := normalize(key)
	for _, ref := range refKeys {
		if normalize(ref) == normKey {
			return ref, "case/underscore mismatch", true
		}
	}
	return "", "", false
}

func normalize(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "-", "_"))
}

// FormatSuggestIssues returns a human-readable summary of suggestions.
func FormatSuggestIssues(issues []SuggestIssue) string {
	if len(issues) == 0 {
		return "no rename suggestions"
	}
	var sb strings.Builder
	for _, iss := range issues {
		sb.WriteString("  " + iss.String() + "\n")
	}
	return sb.String()
}
