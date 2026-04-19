package diff

import "fmt"

// FlatEntry represents a single key-value pair with its source environment name.
type FlatEntry struct {
	Env   string
	Key   string
	Value string
}

// FlattenEnvs takes a map of env name -> key/value pairs and returns a flat list of entries.
func FlattenEnvs(envs map[string]map[string]string) []FlatEntry {
	var entries []FlatEntry
	for env, vars := range envs {
		for k, v := range vars {
			entries = append(entries, FlatEntry{Env: env, Key: k, Value: v})
		}
	}
	return entries
}

// FormatFlatEntries returns a human-readable string for a list of FlatEntry values.
func FormatFlatEntries(entries []FlatEntry) string {
	if len(entries) == 0 {
		return "(no entries)"
	}
	out := ""
	for _, e := range entries {
		out += fmt.Sprintf("[%s] %s=%s\n", e.Env, e.Key, e.Value)
	}
	return out
}

// GroupFlatByKey groups FlatEntry values by key, returning a map of key -> list of FlatEntry.
func GroupFlatByKey(entries []FlatEntry) map[string][]FlatEntry {
	result := make(map[string][]FlatEntry)
	for _, e := range entries {
		result[e.Key] = append(result[e.Key], e)
	}
	return result
}
