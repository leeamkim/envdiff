package diff

import (
	"testing"
)

func baseStatsResult() Result {
	return Result{
		MissingInA: map[string]string{"KEY_A": "val"},
		MissingInB: map[string]string{"KEY_B": "val", "KEY_C": "val"},
		Mismatched: map[string][2]string{"KEY_D": {"x", "y"}},
		Matching:   map[string]string{"KEY_E": "same", "KEY_F": "same"},
	}
}

func TestCompute_Counts(t *testing.T) {
	r := baseStatsResult()
	s := Compute(r)

	if s.TotalKeys != 6 {
		t.Errorf("TotalKeys: want 6, got %d", s.TotalKeys)
	}
	if s.MissingInA != 1 {
		t.Errorf("MissingInA: want 1, got %d", s.MissingInA)
	}
	if s.MissingInB != 2 {
		t.Errorf("MissingInB: want 2, got %d", s.MissingInB)
	}
	if s.Mismatched != 1 {
		t.Errorf("Mismatched: want 1, got %d", s.Mismatched)
	}
	if s.MatchingKeys != 2 {
		t.Errorf("MatchingKeys: want 2, got %d", s.MatchingKeys)
	}
}

func TestCompute_Empty(t *testing.T) {
	s := Compute(Result{
		MissingInA: map[string]string{},
		MissingInB: map[string]string{},
		Mismatched: map[string][2]string{},
		Matching:   map[string]string{},
	})
	if s.TotalKeys != 0 {
		t.Errorf("expected 0 total keys, got %d", s.TotalKeys)
	}
	if s.HasDiff() {
		t.Error("expected HasDiff false for empty result")
	}
}

func TestStats_HasDiff(t *testing.T) {
	s := Compute(baseStatsResult())
	if !s.HasDiff() {
		t.Error("expected HasDiff true")
	}
}

func TestStats_String(t *testing.T) {
	s := Compute(baseStatsResult())
	out := s.String()
	if out == "" {
		t.Error("expected non-empty string")
	}
	expected := "total=6 matching=2 missingInA=1 missingInB=2 mismatched=1"
	if out != expected {
		t.Errorf("String: want %q, got %q", expected, out)
	}
}
