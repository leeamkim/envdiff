package diff

import (
	"fmt"
	"strings"
)

// FormatCoverageReport returns a human-readable coverage report.
func FormatCoverageReport(name string, r CoverageResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Coverage report for %s\n", name)
	fmt.Fprintf(&sb, "  Grade   : %s\n", r.Grade)
	fmt.Fprintf(&sb, "  Coverage: %.1f%% (%d / %d keys)\n",
		r.Coverage*100, r.PresentKeys, r.TotalKeys)
	if len(r.MissingKeys) > 0 {
		fmt.Fprintf(&sb, "  Missing (%d):\n", len(r.MissingKeys))
		for _, k := range r.MissingKeys {
			fmt.Fprintf(&sb, "    - %s\n", k)
		}
	} else {
		fmt.Fprintf(&sb, "  Missing : none\n")
	}
	return sb.String()
}

// FormatMultiCoverageReport renders coverage for multiple environments.
func FormatMultiCoverageReport(results map[string]CoverageResult) string {
	var sb strings.Builder
	names := make([]string, 0, len(results))
	for n := range results {
		names = append(names, n)
	}
	sortStrings(names)

	for _, name := range names {
		sb.WriteString(FormatCoverageReport(name, results[name]))
		sb.WriteString("\n")
	}
	return sb.String()
}
