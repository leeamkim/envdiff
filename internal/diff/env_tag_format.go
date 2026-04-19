package diff

import (
	"fmt"
	"strings"
)

// TagSummary holds aggregated tag counts.
type TagSummary struct {
	Counts map[string]int
	Total  int
}

// SummarizeTags aggregates tag counts from a list of TagEntry.
func SummarizeTags(entries []TagEntry) TagSummary {
	counts := map[string]int{}
	for _, e := range entries {
		for _, tag := range e.Tags {
			counts[tag]++
		}
	}
	total := 0
	for _, v := range counts {
		total += v
	}
	return TagSummary{Counts: counts, Total: total}
}

// FormatTagSummary returns a human-readable summary of tag counts.
func FormatTagSummary(s TagSummary) string {
	if s.Total == 0 {
		return "no tags found"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("total tagged: %d\n", s.Total))
	for tag, count := range s.Counts {
		sb.WriteString(fmt.Sprintf("  %s: %d\n", tag, count))
	}
	return sb.String()
}

// TagEnvMulti applies tag rules across multiple named envs and returns per-env results.
func TagEnvMulti(envs map[string]map[string]string, rules ...TagRule) map[string][]TagEntry {
	result := map[string][]TagEntry{}
	for name, env := range envs {
		result[name] = TagEnv(env, rules...)
	}
	return result
}
