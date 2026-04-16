package diff

import "fmt"

// Summary holds counts of diff results.
type Summary struct {
	MissingInB  int
	MissingInA  int
	Mismatched  int
	Total       int
}

// Summarize returns a Summary from a Result slice.
func Summarize(results []Result) Summary {
	s := Summary{}
	for _, r := range results {
		switch r.Status {
		case MissingInB:
			s.MissingInB++
		case MissingInA:
			s.MissingInA++
		case Mismatched:
			s.Mismatched++
		}
	}
	s.Total = s.MissingInA + s.MissingInB + s.Mismatched
	return s
}

// String returns a human-readable one-line summary.
func (s Summary) String() string {
	if s.Total == 0 {
		return "No differences found."
	}
	return fmt.Sprintf(
		"Summary: %d difference(s) — %d missing in B, %d missing in A, %d mismatched",
		s.Total, s.MissingInB, s.MissingInA, s.Mismatched,
	)
}

// HasDiff returns true when there is at least one difference.
func (s Summary) HasDiff() bool {
	return s.Total > 0
}
