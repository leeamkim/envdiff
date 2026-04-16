package diff

import (
	"testing"
)

func TestCompare_NoDiff(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1", "B": "2"}
	r := Compare(left, right)
	if r.HasDiff() {
		t.Error("expected no diff")
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"A": "1"}
	r := Compare(left, right)
	if len(r.MissingInRight) != 1 || r.MissingInRight[0] != "B" {
		t.Errorf("expected B missing in right, got %v", r.MissingInRight)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := map[string]string{"A": "1"}
	right := map[string]string{"A": "1", "C": "3"}
	r := Compare(left, right)
	if len(r.MissingInLeft) != 1 || r.MissingInLeft[0] != "C" {
		t.Errorf("expected C missing in left, got %v", r.MissingInLeft)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	left := map[string]string{"A": "1", "B": "old"}
	right := map[string]string{"A": "1", "B": "new"}
	r := Compare(left, right)
	if len(r.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(r.Mismatched))
	}
	if v, ok := r.Mismatched["B"]; !ok || v[0] != "old" || v[1] != "new" {
		t.Errorf("unexpected mismatch value: %v", r.Mismatched["B"])
	}
}

func TestCompare_HasDiff(t *testing.T) {
	left := map[string]string{"X": "1"}
	right := map[string]string{"Y": "2"}
	r := Compare(left, right)
	if !r.HasDiff() {
		t.Error("expected diff")
	}
}
