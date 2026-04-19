package diff

import (
	"testing"
)

func TestPromote_Basic(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dest := map[string]string{}
	r := Promote(src, dest)
	if len(r.Added) != 2 {
		t.Fatalf("expected 2 added, got %d", len(r.Added))
	}
	if dest["A"] != "1" || dest["B"] != "2" {
		t.Error("dest not updated correctly")
	}
}

func TestPromote_Conflict_NoOverwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dest := map[string]string{"A": "old"}
	r := Promote(src, dest)
	if len(r.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(r.Conflicts))
	}
	if dest["A"] != "old" {
		t.Error("dest should not have been modified")
	}
}

func TestPromote_Conflict_Overwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dest := map[string]string{"A": "old"}
	r := Promote(src, dest, PromoteOverwrite())
	if len(r.Added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(r.Added))
	}
	if dest["A"] != "new" {
		t.Error("dest should have been overwritten")
	}
}

func TestPromote_SpecificKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dest := map[string]string{}
	r := Promote(src, dest, PromoteKeys("A", "C"))
	if len(r.Added) != 2 {
		t.Fatalf("expected 2 added, got %d", len(r.Added))
	}
	if _, ok := dest["B"]; ok {
		t.Error("B should not have been promoted")
	}
}

func TestPromote_MissingKeyInSrc(t *testing.T) {
	src := map[string]string{"A": "1"}
	dest := map[string]string{}
	r := Promote(src, dest, PromoteKeys("A", "Z"))
	if len(r.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(r.Skipped))
	}
}

func TestFormatPromoteResult(t *testing.T) {
	r := PromoteResult{
		Added:     map[string]string{"A": "1"},
		Conflicts: map[string]string{"B": "2"},
		Skipped:   map[string]string{},
	}
	out := FormatPromoteResult(r)
	if out != "promoted=1 conflicts=1 skipped=0" {
		t.Errorf("unexpected format: %s", out)
	}
}
