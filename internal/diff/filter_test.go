package diff

import (
	"testing"
)

func baseResult() Result {
	return Result{
		MissingInB: map[string]string{"FOO": "bar"},
		MissingInA: map[string]string{"BAZ": "qux"},
		Mismatched: map[string][2]string{"PORT": {"8080", "9090"}},
	}
}

func TestFilter_NoOptions(t *testing.T) {
	r := baseResult()
	out := Filter(r, FilterOptions{})
	if len(out.MissingInB) != 1 || len(out.MissingInA) != 1 || len(out.Mismatched) != 1 {
		t.Error("expected full result when no filter options set")
	}
}

func TestFilter_OnlyMissing(t *testing.T) {
	r := baseResult()
	out := Filter(r, FilterOptions{OnlyMissing: true})
	if len(out.MissingInB) != 1 {
		t.Errorf("expected 1 MissingInB, got %d", len(out.MissingInB))
	}
	if len(out.MissingInA) != 1 {
		t.Errorf("expected 1 MissingInA, got %d", len(out.MissingInA))
	}
	if len(out.Mismatched) != 0 {
		t.Errorf("expected 0 Mismatched, got %d", len(out.Mismatched))
	}
}

func TestFilter_OnlyMismatched(t *testing.T) {
	r := baseResult()
	out := Filter(r, FilterOptions{OnlyMismatched: true})
	if len(out.Mismatched) != 1 {
		t.Errorf("expected 1 Mismatched, got %d", len(out.Mismatched))
	}
	if len(out.MissingInB) != 0 {
		t.Errorf("expected 0 MissingInB, got %d", len(out.MissingInB))
	}
	if len(out.MissingInA) != 0 {
		t.Errorf("expected 0 MissingInA, got %d", len(out.MissingInA))
	}
}

func TestFilter_BothOptions(t *testing.T) {
	r := baseResult()
	out := Filter(r, FilterOptions{OnlyMissing: true, OnlyMismatched: true})
	if len(out.MissingInB) != 1 || len(out.MissingInA) != 1 || len(out.Mismatched) != 1 {
		t.Error("expected full result when both filter options set")
	}
}
