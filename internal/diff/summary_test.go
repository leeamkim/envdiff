package diff

import (
	"testing"
)

func TestSummarize_Empty(t *testing.T) {
	s := Summarize([]Result{})
	if s.Total != 0 {
		t.Errorf("expected 0 total, got %d", s.Total)
	}
	if s.HasDiff() {
		t.Error("expected HasDiff() == false")
	}
	if s.String() != "No differences found." {
		t.Errorf("unexpected string: %s", s.String())
	}
}

func TestSummarize_Counts(t *testing.T) {
	results := []Result{
		{Key: "A", Status: MissingInB},
		{Key: "B", Status: MissingInA},
		{Key: "C", Status: MissingInA},
		{Key: "D", Status: Mismatched},
	}
	s := Summarize(results)
	if s.MissingInB != 1 {
		t.Errorf("expected MissingInB=1, got %d", s.MissingInB)
	}
	if s.MissingInA != 2 {
		t.Errorf("expected MissingInA=2, got %d", s.MissingInA)
	}
	if s.Mismatched != 1 {
		t.Errorf("expected Mismatched=1, got %d", s.Mismatched)
	}
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if !s.HasDiff() {
		t.Error("expected HasDiff() == true")
	}
}

func TestSummarize_String(t *testing.T) {
	results := []Result{
		{Key: "X", Status: Mismatched},
	}
	s := Summarize(results)
	expected := "Summary: 1 difference(s) — 0 missing in B, 0 missing in A, 1 mismatched"
	if s.String() != expected {
		t.Errorf("unexpected string:\ngot:  %s\nwant: %s", s.String(), expected)
	}
}
