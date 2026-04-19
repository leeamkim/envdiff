package diff

import "fmt"

// EnvScore represents a scored comparison between two env files.
type EnvScore struct {
	FileA      string
	FileB      string
	Total      int
	Matched    int
	Missing    int
	Mismatched int
	Score      float64
	Grade      string
}

func (s EnvScore) String() string {
	return fmt.Sprintf("%s vs %s: %.1f%% (%s) [matched=%d missing=%d mismatched=%d]",
		s.FileA, s.FileB, s.Score, s.Grade, s.Matched, s.Missing, s.Mismatched)
}

// ScoreEnvs computes a similarity score between two env maps.
func ScoreEnvs(nameA, nameB string, a, b map[string]string) EnvScore {
	keys := make(map[string]struct{})
	for k := range a {
		keys[k] = struct{}{}
	}
	for k := range b {
		keys[k] = struct{}{}
	}

	total := len(keys)
	if total == 0 {
		return EnvScore{FileA: nameA, FileB: nameB, Grade: "A", Score: 100.0}
	}

	matched, missing, mismatched := 0, 0, 0
	for k := range keys {
		va, okA := a[k]
		vb, okB := b[k]
		switch {
		case okA && okB && va == vb:
			matched++
		case okA && okB:
			mismatched++
		default:
			missing++
		}
	}

	score := float64(matched) / float64(total) * 100.0
	grade := envScoreGrade(score)

	return EnvScore{
		FileA:      nameA,
		FileB:      nameB,
		Total:      total,
		Matched:    matched,
		Missing:    missing,
		Mismatched: mismatched,
		Score:      score,
		Grade:      grade,
	}
}

func envScoreGrade(score float64) string {
	switch {
	case score >= 95:
		return "A"
	case score >= 80:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}
