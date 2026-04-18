package diff

import (
	"testing"
)

func TestNewBaseline_CopiesEntries(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	b := NewBaseline(env)
	if b.Entries["A"] != "1" || b.Entries["B"] != "2" {
		t.Fatal("expected entries to match")
	}
	env["A"] = "mutated"
	if b.Entries["A"] == "mutated" {
		t.Fatal("baseline should not be mutated")
	}
}

func TestMarshalUnmarshalBaseline(t *testing.T) {
	b := NewBaseline(map[string]string{"KEY": "val"})
	data, err := MarshalBaseline(b)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	b2, err := UnmarshalBaseline(data)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if b2.Entries["KEY"] != "val" {
		t.Fatalf("expected KEY=val, got %q", b2.Entries["KEY"])
	}
}

func TestCompareToBaseline_Added(t *testing.T) {
	baseline := NewBaseline(map[string]string{"A": "1"})
	current := map[string]string{"A": "1", "B": "2"}
	d := CompareToBaseline(baseline, current)
	if _, ok := d.Added["B"]; !ok {
		t.Fatal("expected B to be added")
	}
	if d.HasDiff() == false {
		t.Fatal("expected HasDiff true")
	}
}

func TestCompareToBaseline_Removed(t *testing.T) {
	baseline := NewBaseline(map[string]string{"A": "1", "B": "2"})
	current := map[string]string{"A": "1"}
	d := CompareToBaseline(baseline, current)
	if _, ok := d.Removed["B"]; !ok {
		t.Fatal("expected B to be removed")
	}
}

func TestCompareToBaseline_Changed(t *testing.T) {
	baseline := NewBaseline(map[string]string{"A": "old"})
	current := map[string]string{"A": "new"}
	d := CompareToBaseline(baseline, current)
	if pair, ok := d.Changed["A"]; !ok || pair[0] != "old" || pair[1] != "new" {
		t.Fatal("expected A to be changed from old to new")
	}
}

func TestCompareToBaseline_NoDiff(t *testing.T) {
	baseline := NewBaseline(map[string]string{"A": "1"})
	current := map[string]string{"A": "1"}
	d := CompareToBaseline(baseline, current)
	if d.HasDiff() {
		t.Fatal("expected no diff")
	}
}
