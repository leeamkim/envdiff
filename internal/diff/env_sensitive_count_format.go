package diff

import (
	"fmt"
	"sort"
	"strings"
)

// SensitiveCountSummary holds aggregated totals across multiple environments.
type SensitiveCountSummary struct {
	TotalEnvs        int
	TotalKeys        int
	TotalSensitive   int
	TotalRedacted    int
	PerEnv           []SensitiveCountReport
}

// SummarizeSensitiveCounts aggregates SensitiveCountReport results into a
// single summary suitable for multi-environment reporting.
func SummarizeSensitiveCounts(reports []SensitiveCountReport) SensitiveCountSummary {
	sorted := make([]SensitiveCountReport, len(reports))
	copy(sorted, reports)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].EnvName < sorted[j].EnvName
	})

	var totalKeys, totalSensitive, totalRedacted int
	for _, r := range sorted {
		totalKeys += r.TotalKeys
		totalSensitive += r.SensitiveKeys
		totalRedacted += r.RedactedKeys
	}

	return SensitiveCountSummary{
		TotalEnvs:      len(sorted),
		TotalKeys:      totalKeys,
		TotalSensitive: totalSensitive,
		TotalRedacted:  totalRedacted,
		PerEnv:         sorted,
	}
}

// FormatSensitiveCountSummary renders a multi-environment sensitive-key summary
// as a human-readable table.
func FormatSensitiveCountSummary(summary SensitiveCountSummary) string {
	if summary.TotalEnvs == 0 {
		return "no environments to summarize\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sensitive Key Summary (%d environment(s))\n", summary.TotalEnvs))
	sb.WriteString(strings.Repeat("-", 56) + "\n")
	sb.WriteString(fmt.Sprintf("%-20s %8s %9s %9s\n", "Environment", "Total", "Sensitive", "Redacted"))
	sb.WriteString(strings.Repeat("-", 56) + "\n")

	for _, r := range summary.PerEnv {
		sb.WriteString(fmt.Sprintf("%-20s %8d %9d %9d\n",
			r.EnvName, r.TotalKeys, r.SensitiveKeys, r.RedactedKeys))
	}

	sb.WriteString(strings.Repeat("-", 56) + "\n")
	sb.WriteString(fmt.Sprintf("%-20s %8d %9d %9d\n",
		"TOTAL", summary.TotalKeys, summary.TotalSensitive, summary.TotalRedacted))

	return sb.String()
}

// FormatSensitiveCountCSV renders the summary as CSV for machine consumption.
func FormatSensitiveCountCSV(summary SensitiveCountSummary) string {
	var sb strings.Builder
	sb.WriteString("env,total_keys,sensitive_keys,redacted_keys\n")
	for _, r := range summary.PerEnv {
		sb.WriteString(fmt.Sprintf("%s,%d,%d,%d\n",
			r.EnvName, r.TotalKeys, r.SensitiveKeys, r.RedactedKeys))
	}
	return sb.String()
}
