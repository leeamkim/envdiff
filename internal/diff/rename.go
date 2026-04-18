package diff

import "fmt"

// RenameMap maps old key names to new key names.
type RenameMap map[string]string

// RenameEntry records a single rename operation applied to a result.
type RenameEntry struct {
	OldKey string
	NewKey string
}

// ApplyRenames rewrites keys in the given env map according to the rename map,
// returning the updated map and a list of renames that were applied.
func ApplyRenames(env map[string]string, renames RenameMap) (map[string]string, []RenameEntry) {
	updated := make(map[string]string, len(env))
	var applied []RenameEntry

	renamed := make(map[string]bool)
	for oldKey, newKey := range renames {
		if val, ok := env[oldKey]; ok {
			updated[newKey] = val
			applied = append(applied, RenameEntry{OldKey: oldKey, NewKey: newKey})
			renamed[oldKey] = true
		}
	}

	for k, v := range env {
		if !renamed[k] {
			updated[k] = v
		}
	}

	return updated, applied
}

// ParseRenameMap parses a slice of "OLD=NEW" strings into a RenameMap.
func ParseRenameMap(pairs []string) (RenameMap, error) {
	rm := make(RenameMap, len(pairs))
	for _, p := range pairs {
		var oldKey, newKey string
		if _, err := fmt.Sscanf(p, "%s", &oldKey); err != nil {
			return nil, fmt.Errorf("invalid rename pair %q", p)
		}
		found := false
		for i, c := range p {
			if c == '=' {
				oldKey = p[:i]
				newKey = p[i+1:]
				found = true
				break
			}
		}
		if !found || oldKey == "" || newKey == "" {
			return nil, fmt.Errorf("invalid rename pair %q: expected OLD=NEW", p)
		}
		rm[oldKey] = newKey
	}
	return rm, nil
}
