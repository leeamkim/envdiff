package diff

import (
	"testing"
)

var baseSortResult = Result{
	MissingInA: []string{"ZEBRA"},
	MissingInB: []string{"ALPHA"},
	Mismatched: map[string][2]string{
		"MIDDLE": {"v1", "v2"},
	},
}

func TestFlatten_ReturnsAllEntries(t *testing.T) {
	entries := Flatten(baseSortResult)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestSortEntries_ByKey(t *testing.T) {
	entries := Flatten(baseSortResult)
	sorted := SortEntries(entries, SortByKey)
	keys := make([]string, len(sorted))
	for i, e := range sorted {
		keys[i] = e.Key
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("expected ascending key order, got %v", keys)
		}
	}
}

func TestSortEntries_ByKeyDesc(t *testing.T) {
	entries := Flatten(baseSortResult)
	sorted := SortEntries(entries, SortByKeyDesc)
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Key > sorted[i-1].Key {
			t.Errorf("expected descending key order")
		}
	}
}

func TestSortEntries_ByStatus(t *testing.T) {
	entries := Flatten(baseSortResult)
	sorted := SortEntries(entries, SortByStatus)
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Status < sorted[i-1].Status {
			t.Errorf("expected status-grouped order")
		}
	}
}

func TestSortEntries_DefaultIsKey(t *testing.T) {
	entries := Flatten(baseSortResult)
	sorted := SortEntries(entries, "")
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Key < sorted[i-1].Key {
			t.Errorf("expected default key sort")
		}
	}
}
