package diff

import "encoding/json"

// Baseline represents a saved snapshot of a diff result for future comparison.
type Baseline struct {
	Entries map[string]string `json:"entries"`
}

// NewBaseline creates a Baseline from a flat key-value map.
func NewBaseline(env map[string]string) Baseline {
	entries := make(map[string]string, len(env))
	for k, v := range env {
		entries[k] = v
	}
	return Baseline{Entries: entries}
}

// MarshalBaseline serializes a Baseline to JSON bytes.
func MarshalBaseline(b Baseline) ([]byte, error) {
	return json.MarshalIndent(b, "", "  ")
}

// UnmarshalBaseline deserializes a Baseline from JSON bytes.
func UnmarshalBaseline(data []byte) (Baseline, error) {
	var b Baseline
	err := json.Unmarshal(data, &b)
	return b, err
}

// DiffFromBaseline compares a current env map against a saved baseline.
// Returns keys added, removed, or changed relative to the baseline.
type BaselineDiff struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [baseline, current]
}

func CompareToBaseline(baseline Baseline, current map[string]string) BaselineDiff {
	result := BaselineDiff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}
	for k, v := range current {
		if bv, ok := baseline.Entries[k]; !ok {
			result.Added[k] = v
		} else if bv != v {
			result.Changed[k] = [2]string{bv, v}
		}
	}
	for k, v := range baseline.Entries {
		if _, ok := current[k]; !ok {
			result.Removed[k] = v
		}
	}
	return result
}

// HasDiff returns true if there are any differences from the baseline.
func (d BaselineDiff) HasDiff() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
