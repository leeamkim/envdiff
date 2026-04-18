package diff

import (
	"strings"
	"testing"
)

func baseScorecardResult() Result {
	return Result{
		Same:       map[string]string{"HOST": "localhost", "PORT": "5432"},
		MissingInB: map[string]string{"SECRET": "abc"},
		MissingInA: map[string]string{"NEW_KEY": "xyz"},
		Mismatched: map[string][2]string{"DB": {"old", "new"}},
	}
}

func TestComputeScorecard_Counts(t *testing.T) {
	sc := ComputeScorecard(baseScorecardResult())
	if sc.Total != 5 {
		t.Errorf("expected total 5, got %d", sc.Total)
	}
	if sc.Matched != 2 {
		t.Errorf("expected matched 2, got %d", sc.Matched)
	}
	if sc.Missing != 2 {
		t.Errorf("expected missing 2, got %d", sc.Missing)
	}
	if sc.Mismatched != 1 {
		t.Errorf("expected mismatched 1, got %d", sc.Mismatched)
	}
}

func TestComputeScorecard_Score(t *testing.T) {
	sc := ComputeScorecard(baseScorecardResult())
	expected := 40.0
	if sc.Score != expected {
		t.Errorf("expected score %.1f, got %.1f", expected, sc.Score)
	}
}

func TestComputeScorecard_Empty(t *testing.T) {
	sc := ComputeScorecard(Result{})
	if sc.Score != 0 {
		t.Errorf("expected score 0, got %.1f", sc.Score)
	}
	if sc.Total != 0 {
		t.Errorf("expected total 0, got %d", sc.Total)
	}
}

func TestScorecard_Grade(t *testing.T) {
	cases := []struct{ score float64; grade string }{
		{100, "A"}, {90, "A"}, {80, "B"}, {75, "B"},
		{65, "C"}, {50, "D"}, {30, "F"},
	}
	for _, c := range cases {
		sc := Scorecard{Score: c.score}
		if sc.Grade() != c.grade {
			t.Errorf("score %.0f: expected grade %s, got %s", c.score, c.grade, sc.Grade())
		}
	}
}

func TestScorecard_String(t *testing.T) {
	sc := Scorecard{Total: 5, Matched: 2, Missing: 2, Mismatched: 1, Score: 40.0}
	s := sc.String()
	if !strings.Contains(s, "40.0%") {
		t.Errorf("expected score in string, got: %s", s)
	}
	if !strings.Contains(s, "F") {
		t.Errorf("expected grade F in string, got: %s", s)
	}
}
