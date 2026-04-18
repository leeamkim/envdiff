package diff

import (
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}
	r := Merge(a, b, false)
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", r.Conflicts)
	}
	if r.Merged["FOO"] != "1" || r.Merged["BAR"] != "2" || r.Merged["BAZ"] != "3" {
		t.Errorf("unexpected merged map: %v", r.Merged)
	}
}

func TestMerge_PreferB(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	r := Merge(a, b, false)
	if r.Merged["FOO"] != "new" {
		t.Errorf("expected 'new', got %s", r.Merged["FOO"])
	}
	if len(r.Conflicts) != 1 || r.Conflicts[0] != "FOO" {
		t.Errorf("expected conflict on FOO, got %v", r.Conflicts)
	}
}

func TestMerge_PreferA(t *testing.T) {
	a := map[string]string{"FOO": "old"}
	b := map[string]string{"FOO": "new"}
	r := Merge(a, b, true)
	if r.Merged["FOO"] != "old" {
		t.Errorf("expected 'old', got %s", r.Merged["FOO"])
	}
	if len(r.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %v", r.Conflicts)
	}
}

func TestMerge_MultipleConflictsSorted(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "1", "M": "1"}
	b := map[string]string{"Z": "2", "A": "2", "M": "2"}
	r := Merge(a, b, false)
	if len(r.Conflicts) != 3 {
		t.Fatalf("expected 3 conflicts, got %d", len(r.Conflicts))
	}
	if r.Conflicts[0] != "A" || r.Conflicts[1] != "M" || r.Conflicts[2] != "Z" {
		t.Errorf("conflicts not sorted: %v", r.Conflicts)
	}
}
