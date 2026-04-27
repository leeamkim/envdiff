package diff

import (
	"fmt"
	"strings"
)

// ObsoleteSummary holds aggregate statistics for obsolete key analysis.
type ObsoleteSummary struct {
	Total    int
	Obsolete int
	Clean    bool
}

func (s ObsoleteSummary) String() string {
	if s.Clean {
		return fmt.Sprintf("all %d keys are current (no obsolete keys)", s.Total)
	}
	return fmt.Sprintf("%d/%d keys are obsolete", s.Obsolete, s.Total)
}

// SummarizeObsolete computes summary statistics for an obsolete check result.
func SummarizeObsolete(env map[string]string, issues []ObsoleteIssue) ObsoleteSummary {
	s := ObsoleteSummary{
		Total:    len(env),
		Obsolete: len(issues),
	}
	s.Clean = s.Obsolete == 0
	return s
}

// FormatObsoleteReport renders a full report including summary and per-key details.
func FormatObsoleteReport(env map[string]string, issues []ObsoleteIssue) string {
	summary := SummarizeObsolete(env, issues)
	var sb strings.Builder
	sb.WriteString(summary.String() + "\n")
	if len(issues) > 0 {
		sb.WriteString("obsolete keys:\n")
		for _, issue := range issues {
			sb.WriteString(fmt.Sprintf("  %s\n", issue.String()))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
