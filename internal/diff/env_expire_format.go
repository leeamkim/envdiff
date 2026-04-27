package diff

import (
	"fmt"
	"strings"
)

// ExpireSummary holds aggregate counts from a set of ExpireIssues.
type ExpireSummary struct {
	Total   int
	Expired int
	Warning int
}

// SummarizeExpiry aggregates expire issues into a summary.
func SummarizeExpiry(issues []ExpireIssue) ExpireSummary {
	s := ExpireSummary{Total: len(issues)}
	for _, issue := range issues {
		if issue.Expired {
			s.Expired++
		} else {
			s.Warning++
		}
	}
	return s
}

func (s ExpireSummary) String() string {
	return fmt.Sprintf("total=%d expired=%d warning=%d", s.Total, s.Expired, s.Warning)
}

// FormatExpireReport returns a detailed multi-section report.
func FormatExpireReport(issues []ExpireIssue) string {
	if len(issues) == 0 {
		return "no expiry issues found"
	}

	var sb strings.Builder
	summary := SummarizeExpiry(issues)
	sb.WriteString(fmt.Sprintf("expiry report: %s\n", summary))
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	for _, issue := range issues {
		label := "WARNING"
		if issue.Expired {
			label = "EXPIRED"
		}
		sb.WriteString(fmt.Sprintf("[%s] %s\n", label, issue))
	}

	return strings.TrimRight(sb.String(), "\n")
}
