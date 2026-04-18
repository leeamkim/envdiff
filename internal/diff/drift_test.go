package diff

import (
	"strings"
	"testing"
)

func TestDetectDrift_NoDrift(t *testing.T) {
	ref := map[string]string{"A": "1", "B": "2"}
	cur := map[string]string{"A": "1", "B": "2"}
	r := DetectDrift(ref, cur)
	if r.HasDrift() {
		t.Errorf("expected no drift, got %d entries", len(r.Entries))
	}
	if r.Summary() != "no drift detected" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestDetectDrift_Added(t *testing.T) {
	ref := map[string]string{"A": "1"}
	cur := map[string]string{"A": "1", "B": "2"}
	r := DetectDrift(ref, cur)
	if !r.HasDrift() {
		t.Fatal("expected drift")
	}
	if r.Entries[0].Status != "added" || r.Entries[0].Key != "B" {
		t.Errorf("unexpected entry: %+v", r.Entries[0])
	}
}

func TestDetectDrift_Removed(t *testing.T) {
	ref := map[string]string{"A": "1", "B": "2"}
	cur := map[string]string{"A": "1"}
	r := DetectDrift(ref, cur)
	if !r.HasDrift() {
		t.Fatal("expected drift")
	}
	if r.Entries[0].Status != "removed" || r.Entries[0].Key != "B" {
		t.Errorf("unexpected entry: %+v", r.Entries[0])
	}
}

func TestDetectDrift_Changed(t *testing.T) {
	ref := map[string]string{"A": "old"}
	cur := map[string]string{"A": "new"}
	r := DetectDrift(ref, cur)
	if !r.HasDrift() {
		t.Fatal("expected drift")
	}
	e := r.Entries[0]
	if e.Status != "changed" || e.OldValue != "old" || e.NewValue != "new" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestDriftEntry_String(t *testing.T) {
	cases := []struct {
		entry  DriftEntry
		expect string
	}{
		{DriftEntry{Key: "X", NewValue: "v", Status: "added"}, "[added]"},
		{DriftEntry{Key: "X", OldValue: "v", Status: "removed"}, "[removed]"},
		{DriftEntry{Key: "X", OldValue: "a", NewValue: "b", Status: "changed"}, "[changed]"},
	}
	for _, c := range cases {
		if !strings.Contains(c.entry.String(), c.expect) {
			t.Errorf("expected %q in %q", c.expect, c.entry.String())
		}
	}
}

func TestDriftReport_Summary(t *testing.T) {
	r := &DriftReport{Entries: []DriftEntry{{Key: "A", Status: "added"}, {Key: "B", Status: "removed"}}}
	if r.Summary() != "2 key(s) drifted" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}
