package diff

import (
	"fmt"
	"strings"
)

// HealthStatus represents the overall health of an env file.
type HealthStatus struct {
	TotalKeys    int
	EmptyValues  int
	Placeholders int
	LowercaseKeys int
	Score        int // 0-100
	Grade        string
}

func (h HealthStatus) String() string {
	return fmt.Sprintf("Grade: %s (score=%d, total=%d, empty=%d, placeholders=%d, lowercase=%d)",
		h.Grade, h.Score, h.TotalKeys, h.EmptyValues, h.Placeholders, h.LowercaseKeys)
}

// ComputeHealth evaluates the quality of an env map and returns a HealthStatus.
// It penalizes empty values (-10 each), placeholder values (-8 each),
// and lowercase keys (-5 each) to produce a score from 0 to 100.
func ComputeHealth(env map[string]string) HealthStatus {
	total := len(env)
	if total == 0 {
		return HealthStatus{Grade: "N/A"}
	}

	var empty, placeholders, lowercase int
	for k, v := range env {
		if v == "" {
			empty++
		}
		if isPlaceholder(v) {
			placeholders++
		}
		if k != strings.ToUpper(k) {
			lowercase++
		}
	}

	penalty := 0
	penalty += (empty * 10)
	penalty += (placeholders * 8)
	penalty += (lowercase * 5)

	score := 100 - penalty
	if score < 0 {
		score = 0
	}

	grade := scoreGrade(score)
	return HealthStatus{
		TotalKeys:     total,
		EmptyValues:   empty,
		Placeholders:  placeholders,
		LowercaseKeys: lowercase,
		Score:         score,
		Grade:         grade,
	}
}

// isPlaceholder reports whether a value looks like an unfilled placeholder.
func isPlaceholder(v string) bool {
	switch strings.ToUpper(v) {
	case "CHANGEME", "TODO", "PLACEHOLDER", "FIXME", "YOUR_VALUE_HERE":
		return true
	}
	return false
}

func scoreGrade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}
