package diff

import (
	"fmt"
	"strings"
)

// DependencySummary holds aggregated counts for a dependency check run.
type DependencySummary struct {
	Total   int
	Missing int
	Empty   int
}

// SummarizeDependencies builds a DependencySummary from a slice of issues.
func SummarizeDependencies(issues []DependencyIssue) DependencySummary {
	s := DependencySummary{Total: len(issues)}
	for _, iss := range issues {
		switch iss.Reason {
		case "dependency key is missing":
			s.Missing++
		case "dependency key is empty":
			s.Empty++
		}
	}
	return s
}

// FormatDependencyReport returns a full report string including a summary line.
func FormatDependencyReport(issues []DependencyIssue) string {
	var sb strings.Builder
	if len(issues) == 0 {
		sb.WriteString("dependency check passed: no issues found\n")
		return sb.String()
	}
	summary := SummarizeDependencies(issues)
	sb.WriteString(fmt.Sprintf("dependency issues: %d total (%d missing, %d empty)\n",
		summary.Total, summary.Missing, summary.Empty))
	for _, iss := range issues {
		sb.WriteString(fmt.Sprintf("  [DEPENDENCY] %s\n", iss.String()))
	}
	return sb.String()
}
