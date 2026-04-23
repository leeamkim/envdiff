package diff

import "fmt"

// CoverageResult holds coverage statistics for a set of env files
// compared against a reference (e.g. a template or canonical env).
type CoverageResult struct {
	TotalKeys   int
	PresentKeys int
	MissingKeys []string
	Coverage    float64 // 0.0 – 1.0
	Grade       string
}

func (c CoverageResult) String() string {
	return fmt.Sprintf("coverage=%.1f%% grade=%s present=%d/%d missing=%d",
		c.Coverage*100, c.Grade, c.PresentKeys, c.TotalKeys, len(c.MissingKeys))
}

// ComputeCoverage calculates how many keys from reference exist (with
// non-empty values) in target.
func ComputeCoverage(reference, target map[string]string) CoverageResult {
	if len(reference) == 0 {
		return CoverageResult{Grade: "A"}
	}

	var missing []string
	present := 0

	for k := range reference {
		v, ok := target[k]
		if ok && v != "" {
			present++
		} else {
			missing = append(missing, k)
		}
	}

	sortStrings(missing)

	coverage := float64(present) / float64(len(reference))
	return CoverageResult{
		TotalKeys:   len(reference),
		PresentKeys: present,
		MissingKeys: missing,
		Coverage:    coverage,
		Grade:       coverageGrade(coverage),
	}
}

func coverageGrade(c float64) string {
	switch {
	case c >= 0.95:
		return "A"
	case c >= 0.80:
		return "B"
	case c >= 0.65:
		return "C"
	case c >= 0.50:
		return "D"
	default:
		return "F"
	}
}
