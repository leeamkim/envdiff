package diff

import "fmt"

// Stats holds aggregate statistics for a diff result.
type Stats struct {
	TotalKeys    int
	MissingInA   int
	MissingInB   int
	Mismatched   int
	MatchingKeys int
}

// Compute derives Stats from a Result.
func Compute(r Result) Stats {
	allKeys := make(map[string]struct{})
	for k := range r.MissingInA {
		allKeys[k] = struct{}{}
	}
	for k := range r.MissingInB {
		allKeys[k] = struct{}{}
	}
	for k := range r.Mismatched {
		allKeys[k] = struct{}{}
	}
	for k := range r.Matching {
		allKeys[k] = struct{}{}
	}

	s := Stats{
		TotalKeys:    len(allKeys),
		MissingInA:   len(r.MissingInA),
		MissingInB:   len(r.MissingInB),
		Mismatched:   len(r.Mismatched),
		MatchingKeys: len(r.Matching),
	}
	return s
}

// String returns a human-readable summary line.
func (s Stats) String() string {
	return fmt.Sprintf(
		"total=%d matching=%d missingInA=%d missingInB=%d mismatched=%d",
		s.TotalKeys, s.MatchingKeys, s.MissingInA, s.MissingInB, s.Mismatched,
	)
}

// HasDiff returns true when any discrepancy exists.
func (s Stats) HasDiff() bool {
	return s.MissingInA > 0 || s.MissingInB > 0 || s.Mismatched > 0
}

// DiffCount returns the total number of discrepant keys across all diff
// categories (missing in A, missing in B, and mismatched values).
func (s Stats) DiffCount() int {
	return s.MissingInA + s.MissingInB + s.Mismatched
}
