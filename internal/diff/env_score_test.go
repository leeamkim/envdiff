package diff

import (
	"strings"
	"testing"
)

func TestScoreEnvs_Perfect(t *testing.T) {
	a := map[string]string{"KEY": "val", "FOO": "bar"}
	b := map[string]string{"KEY": "val", "FOO": "bar"}
	s := ScoreEnvs("a.env", "b.env", a, b)
	if s.Score != 100.0 {
		t.Errorf("expected 100, got %.1f", s.Score)
	}
	if s.Grade != "A" {
		t.Errorf("expected A, got %s", s.Grade)
	}
	if s.Matched != 2 {
		t.Errorf("expected matched=2, got %d", s.Matched)
	}
}

func TestScoreEnvs_Missing(t *testing.T) {
	a := map[string]string{"KEY": "val"}
	b := map[string]string{"KEY": "val", "EXTRA": "x"}
	s := ScoreEnvs("a.env", "b.env", a, b)
	if s.Missing != 1 {
		t.Errorf("expected missing=1, got %d", s.Missing)
	}
	if s.Total != 2 {
		t.Errorf("expected total=2, got %d", s.Total)
	}
}

func TestScoreEnvs_Mismatched(t *testing.T) {
	a := map[string]string{"KEY": "val1"}
	b := map[string]string{"KEY": "val2"}
	s := ScoreEnvs("a.env", "b.env", a, b)
	if s.Mismatched != 1 {
		t.Errorf("expected mismatched=1, got %d", s.Mismatched)
	}
	if s.Score != 0.0 {
		t.Errorf("expected score=0, got %.1f", s.Score)
	}
	if s.Grade != "F" {
		t.Errorf("expected F, got %s", s.Grade)
	}
}

func TestScoreEnvs_Empty(t *testing.T) {
	s := ScoreEnvs("a.env", "b.env", map[string]string{}, map[string]string{})
	if s.Score != 100.0 {
		t.Errorf("expected 100 for empty, got %.1f", s.Score)
	}
}

func TestEnvScore_String(t *testing.T) {
	s := ScoreEnvs("a.env", "b.env",
		map[string]string{"A": "1", "B": "2"},
		map[string]string{"A": "1", "B": "2"},
	)
	out := s.String()
	if !strings.Contains(out, "a.env") || !strings.Contains(out, "b.env") {
		t.Errorf("expected file names in output: %s", out)
	}
	if !strings.Contains(out, "100.0%") {
		t.Errorf("expected score in output: %s", out)
	}
}
