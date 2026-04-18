package diff

import "sort"

// MergeResult holds the merged output of two env maps.
type MergeResult struct {
	Merged  map[string]string
	Conflicts []string
}

// Merge combines two env maps into one.
// Values from b take precedence unless preferA is true.
// Keys with differing values are recorded as conflicts.
func Merge(a, b map[string]string, preferA bool) MergeResult {
	merged := make(map[string]string)
	conflicts := []string{}

	for k, v := range a {
		merged[k] = v
	}

	for k, vb := range b {
		va, exists := merged[k]
		if !exists {
			merged[k] = vb
			continue
		}
		if va != vb {
			conflicts = append(conflicts, k)
			if !preferA {
				merged[k] = vb
			}
		}
	}

	sort.Strings(conflicts)
	return MergeResult{Merged: merged, Conflicts: conflicts}
}
