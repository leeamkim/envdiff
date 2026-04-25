package diff

import "sort"

// RenameCandidate represents a potential key rename between two env maps.
type RenameCandidate struct {
	OldKey string
	NewKey string
	Value  string
}

// String returns a human-readable description of the rename candidate.
func (r RenameCandidate) String() string {
	return "rename: " + r.OldKey + " -> " + r.NewKey + " (value: " + r.Value + ")"
}

// DetectRenames finds keys that exist only in one env map but share the same
// value with a key that exists only in the other, suggesting a rename occurred.
func DetectRenames(a, b map[string]string) []RenameCandidate {
	onlyInA := make(map[string]string)
	onlyInB := make(map[string]string)

	for k, v := range a {
		if _, ok := b[k]; !ok {
			onlyInA[k] = v
		}
	}
	for k, v := range b {
		if _, ok := a[k]; !ok {
			onlyInB[k] = v
		}
	}

	// Build reverse map: value -> keys in B
	valToB := make(map[string][]string)
	for k, v := range onlyInB {
		if v != "" {
			valToB[v] = append(valToB[v], k)
		}
	}

	var candidates []RenameCandidate
	for oldKey, val := range onlyInA {
		if val == "" {
			continue
		}
		matches, ok := valToB[val]
		if !ok || len(matches) != 1 {
			continue
		}
		candidates = append(candidates, RenameCandidate{
			OldKey: oldKey,
			NewKey: matches[0],
			Value:  val,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].OldKey < candidates[j].OldKey
	})
	return candidates
}

// FormatRenameReport formats detected rename candidates as a human-readable string.
func FormatRenameReport(candidates []RenameCandidate) string {
	if len(candidates) == 0 {
		return "no rename candidates detected\n"
	}
	out := ""
	for _, c := range candidates {
		out += c.String() + "\n"
	}
	return out
}
