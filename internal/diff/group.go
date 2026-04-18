package diff

import "sort"

// Group represents diff entries grouped by a prefix (e.g., "DB", "AWS").
type Group struct {
	Prefix  string
	Entries []FlatEntry
}

// GroupByPrefix splits flat entries into groups based on key prefix.
// Keys are split on the first "_" character. Keys with no prefix go into "OTHER".
func GroupByPrefix(entries []FlatEntry) []Group {
	groupMap := make(map[string][]FlatEntry)

	for _, e := range entries {
		prefix := extractPrefix(e.Key)
		groupMap[prefix] = append(groupMap[prefix], e)
	}

	groups := make([]Group, 0, len(groupMap))
	for prefix, ents := range groupMap {
		groups = append(groups, Group{Prefix: prefix, Entries: ents})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Prefix < groups[j].Prefix
	})

	return groups
}

func extractPrefix(key string) string {
	for i, ch := range key {
		if ch == '_' && i > 0 {
			return key[:i]
		}
	}
	return "OTHER"
}
