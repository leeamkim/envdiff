package diff

import (
	"testing"
)

func baseGroupEntries() []FlatEntry {
	return []FlatEntry{
		{Key: "DB_HOST", Status: "missing_in_b"},
		{Key: "DB_PORT", Status: "ok"},
		{Key: "AWS_KEY", Status: "mismatched"},
		{Key: "AWS_SECRET", Status: "missing_in_a"},
		{Key: "PORT", Status: "ok"},
	}
}

func TestGroupByPrefix_Keys(t *testing.T) {
	groups := GroupByPrefix(baseGroupEntries())
	prefixes := make([]string, len(groups))
	for i, g := range groups {
		prefixes[i] = g.Prefix
	}
	expected := []string{"AWS", "DB", "OTHER"}
	for i, p := range expected {
		if prefixes[i] != p {
			t.Errorf("expected prefix %q at index %d, got %q", p, i, prefixes[i])
		}
	}
}

func TestGroupByPrefix_Counts(t *testing.T) {
	groups := GroupByPrefix(baseGroupEntries())
	counts := map[string]int{}
	for _, g := range groups {
		counts[g.Prefix] = len(g.Entries)
	}
	if counts["DB"] != 2 {
		t.Errorf("expected 2 DB entries, got %d", counts["DB"])
	}
	if counts["AWS"] != 2 {
		t.Errorf("expected 2 AWS entries, got %d", counts["AWS"])
	}
	if counts["OTHER"] != 1 {
		t.Errorf("expected 1 OTHER entry, got %d", counts["OTHER"])
	}
}

func TestGroupByPrefix_Empty(t *testing.T) {
	groups := GroupByPrefix([]FlatEntry{})
	if len(groups) != 0 {
		t.Errorf("expected no groups, got %d", len(groups))
	}
}

func TestGroupByPrefix_NoUnderscore(t *testing.T) {
	entries := []FlatEntry{
		{Key: "HOST", Status: "ok"},
		{Key: "PORT", Status: "ok"},
	}
	groups := GroupByPrefix(entries)
	if len(groups) != 1 || groups[0].Prefix != "OTHER" {
		t.Errorf("expected single OTHER group, got %+v", groups)
	}
}
