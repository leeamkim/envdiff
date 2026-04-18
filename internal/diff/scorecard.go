package diff

import "fmt"

// Scorecard holds a health score summary for a diff result.
type Scorecard struct {
	Total     int
	Matched   int
	Missing   int
	Mismatched int
	Score     float64 // 0.0 - 100.0
}

// Grade returns a letter grade based on the score.
func (s Scorecard) Grade() string {
	switch {
	case s.Score >= 90:
		return "A"
	case s.Score >= 75:
		return "B"
	case s.Score >= 60:
		return "C"
	case s.Score >= 40:
		return "D"
	default:
		return "F"
	}
}

func (s Scorecard) String() string {
	return fmt.Sprintf("Score: %.1f%% (%s) — %d matched, %d missing, %d mismatched of %d total",
		s.Score, s.Grade(), s.Matched, s.Missing, s.Mismatched, s.Total)
}

// ComputeScorecard calculates a health scorecard from a Result.
func ComputeScorecard(r Result) Scorecard {
	matched := 0
	missing := 0
	mismatched := 0

	for _, v := range r.MissingInA {
		_ = v
		missing++
	}
	for _, v := range r.MissingInB {
		_ = v
		missing++
	}
	for _, v := range r.Mismatched {
		_ = v
		mismatched++
	}

	// Matched = keys in both A and B that are not mismatched
	allKeys := map[string]struct{}{}
	for k := range r.MissingInA {
		allKeys[k] = struct{}{}
	}
	for k := range r.MissingInB {
		allKeys[k] = struct{}{}
	}
	for k := range r.Mismatched {
		allKeys[k] = struct{}{}
	}
	for k := range r.Same {
		allKeys[k] = struct{}{}
		matched++
	}

	total := len(allKeys)
	var score float64
	if total > 0 {
		score = float64(matched) / float64(total) * 100.0
	}

	return Scorecard{
		Total:      total,
		Matched:    matched,
		Missing:    missing,
		Mismatched: mismatched,
		Score:      score,
	}
}
