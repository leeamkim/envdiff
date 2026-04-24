package diff

import "sort"

// OverlapEntry describes a key shared between two or more environments
// along with its consistency status.
type OverlapEntry struct {
	Key        string
	Envs       []string
	Consistent bool
	Values     map[string]string
}

// OverlapResult holds the full overlap analysis across a set of environments.
type OverlapResult struct {
	Entries []OverlapEntry
}

// ComputeOverlap finds all keys that appear in at least minEnvs environments
// and reports whether their values are consistent.
func ComputeOverlap(envs map[string]map[string]string, minEnvs int) OverlapResult {
	if minEnvs < 1 {
		minEnvs = 1
	}

	// Count how many envs each key appears in.
	keyEnvs := map[string][]string{}
	for envName, vars := range envs {
		for k := range vars {
			keyEnvs[k] = append(keyEnvs[k], envName)
		}
	}

	var entries []OverlapEntry
	for key, envNames := range keyEnvs {
		if len(envNames) < minEnvs {
			continue
		}
		sort.Strings(envNames)

		values := map[string]string{}
		for _, en := range envNames {
			values[en] = envs[en][key]
		}

		consistent := true
		var first string
		for i, en := range envNames {
			if i == 0 {
				first = values[en]
			} else if values[en] != first {
				consistent = false
				break
			}
		}

		entries = append(entries, OverlapEntry{
			Key:        key,
			Envs:       envNames,
			Consistent: consistent,
			Values:     values,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return OverlapResult{Entries: entries}
}
