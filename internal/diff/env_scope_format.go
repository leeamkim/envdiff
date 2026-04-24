package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ScopeSummary holds aggregated scope violation counts per environment.
type ScopeSummary struct {
	Env        string
	Violations int
}

// SummarizeScopes groups scope issues by environment and returns a sorted summary.
func SummarizeScopes(issues []ScopeIssue) []ScopeSummary {
	counts := make(map[string]int)
	for _, issue := range issues {
		counts[issue.Env]++
	}

	summaries := make([]ScopeSummary, 0, len(counts))
	for env, count := range counts {
		summaries = append(summaries, ScopeSummary{Env: env, Violations: count})
	}
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Env < summaries[j].Env
	})
	return summaries
}

// FormatScopeSummary renders a scope summary table as a string.
func FormatScopeSummary(summaries []ScopeSummary) string {
	if len(summaries) == 0 {
		return "all scopes clean"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s %s\n", "ENV", "VIOLATIONS"))
	sb.WriteString(strings.Repeat("-", 32) + "\n")
	for _, s := range summaries {
		sb.WriteString(fmt.Sprintf("%-20s %d\n", s.Env, s.Violations))
	}
	return strings.TrimRight(sb.String(), "\n")
}
