package diff

import "sort"

// SortOrder defines how diff results should be sorted.
type SortOrder string

const (
	SortByKey      SortOrder = "key"
	SortByStatus   SortOrder = "status"
	SortByKeyDesc  SortOrder = "key-desc"
)

// DiffEntry represents a flat, sortable diff entry.
type DiffEntry struct {
	Key    string
	Status string // "missing_in_a", "missing_in_b", "mismatched"
	ValueA string
	ValueB string
}

// Flatten converts a Result into a slice of DiffEntry.
func Flatten(r Result) []DiffEntry {
	var entries []DiffEntry
	for _, k := range r.MissingInA {
		entries = append(entries, DiffEntry{Key: k, Status: "missing_in_a"})
	}
	for _, k := range r.MissingInB {
		entries = append(entries, DiffEntry{Key: k, Status: "missing_in_b"})
	}
	for k, v := range r.Mismatched {
		entries = append(entries, DiffEntry{Key: k, Status: "mismatched", ValueA: v[0], ValueB: v[1]})
	}
	return entries
}

// SortEntries sorts a slice of DiffEntry by the given order.
func SortEntries(entries []DiffEntry, order SortOrder) []DiffEntry {
	switch order {
	case SortByStatus:
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Status == entries[j].Status {
				return entries[i].Key < entries[j].Key
			}
			return entries[i].Status < entries[j].Status
		})
	case SortByKeyDesc:
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key > entries[j].Key
		})
	default: // SortByKey
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
	}
	return entries
}
