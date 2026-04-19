package diff

import "fmt"

// TagEntry represents a key tagged with metadata across environments.
type TagEntry struct {
	Key  string
	Tags []string
}

// TagRule is a function that returns tags for a given key/value pair.
type TagRule func(key, value string) []string

// TagRuleSecrets tags keys that look like secrets.
func TagRuleSecrets(key, _ string) []string {
	if isSensitive(key) {
		return []string{"secret"}
	}
	return nil
}

// TagRuleEmpty tags keys with empty values.
func TagRuleEmpty(_ string, value string) []string {
	if value == "" {
		return []string{"empty"}
	}
	return nil
}

// TagRulePlaceholder tags keys with placeholder values.
func TagRulePlaceholder(_ string, value string) []string {
	if value == "CHANGEME" || value == "TODO" || value == "PLACEHOLDER" {
		return []string{"placeholder"}
	}
	return nil
}

// TagEnv applies tag rules to an env map and returns a list of TagEntry.
func TagEnv(env map[string]string, rules ...TagRule) []TagEntry {
	keys := sortedKeys(env)
	var entries []TagEntry
	for _, k := range keys {
		var tags []string
		for _, rule := range rules {
			tags = append(tags, rule(k, env[k])...)
		}
		if len(tags) > 0 {
			entries = append(entries, TagEntry{Key: k, Tags: tags})
		}
	}
	return entries
}

// FormatTagEntries returns a human-readable string of tag entries.
func FormatTagEntries(entries []TagEntry) string {
	if len(entries) == 0 {
		return "no tags"
	}
	out := ""
	for _, e := range entries {
		out += fmt.Sprintf("%s: %v\n", e.Key, e.Tags)
	}
	return out
}
