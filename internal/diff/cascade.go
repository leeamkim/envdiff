package diff

import "fmt"

// CascadeEntry represents a key that propagates across a chain of environments.
type CascadeEntry struct {
	Key    string
	Values map[string]string // env name -> value
	Drift  bool             // true if any value differs from the first
}

func (c CascadeEntry) String() string {
	if c.Drift {
		return fmt.Sprintf("%s [DRIFT across %d envs]", c.Key, len(c.Values))
	}
	return fmt.Sprintf("%s [consistent across %d envs]", c.Key, len(c.Values))
}

// CascadeResult holds the full cascade analysis across N environments.
type CascadeResult struct {
	EnvNames []string
	Entries  []CascadeEntry
}

// Cascade compares a key across multiple named environments.
// envs is a map of env-name -> parsed key/value map.
func Cascade(envNames []string, envs map[string]map[string]string) CascadeResult {
	keySet := map[string]struct{}{}
	for _, m := range envs {
		for k := range m {
			keySet[k] = struct{}{}
		}
	}

	var entries []CascadeEntry
	for key := range keySet {
		values := map[string]string{}
		var first string
		firstSet := false
		drift := false
		for _, name := range envNames {
			v, ok := envs[name][key]
			if ok {
				values[name] = v
				if !firstSet {
					first = v
					firstSet = true
				} else if v != first {
					drift = true
				}
			}
		}
		entries = append(entries, CascadeEntry{Key: key, Values: values, Drift: drift})
	}
	SortEntries(Flatten(Result{Mismatched: map[string][2]string{}}))
	return CascadeResult{EnvNames: envNames, Entries: entries}
}
