package diff

import (
	"fmt"
	"sort"
	"strings"
)

// SensitiveCountEntry holds the count of sensitive keys for a single env.
type SensitiveCountEntry struct {
	EnvName        string
	TotalKeys      int
	SensitiveKeys  int
	RedactedKeys   int
}

// SensitiveCountResult holds entries for all envs.
type SensitiveCountResult struct {
	Entries []SensitiveCountEntry
}

var defaultSensitivePatterns = []string{
	"password", "secret", "token", "key", "api_key",
	"auth", "credential", "private", "passwd",
}

// CountSensitiveKeys scans each env map and counts keys that match sensitive
// patterns. redactedKeys is an optional set of keys already known to be redacted.
func CountSensitiveKeys(envs map[string]map[string]string, redactedKeys []string) SensitiveCountResult {
	redacted := make(map[string]struct{}, len(redactedKeys))
	for _, k := range redactedKeys {
		redacted[strings.ToLower(k)] = struct{}{}
	}

	names := make([]string, 0, len(envs))
	for name := range envs {
		names = append(names, name)
	}
	sort.Strings(names)

	var entries []SensitiveCountEntry
	for _, name := range names {
		vars := envs[name]
		entry := SensitiveCountEntry{
			EnvName:   name,
			TotalKeys: len(vars),
		}
		for k := range vars {
			lower := strings.ToLower(k)
			if isSensitiveKey(lower) {
				entry.SensitiveKeys++
				if _, ok := redacted[lower]; ok {
					entry.RedactedKeys++
				}
			}
		}
		entries = append(entries, entry)
	}
	return SensitiveCountResult{Entries: entries}
}

func isSensitiveKey(lower string) bool {
	for _, p := range defaultSensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// FormatSensitiveCountReport returns a human-readable table of sensitive key counts.
func FormatSensitiveCountReport(r SensitiveCountResult) string {
	if len(r.Entries) == 0 {
		return "no environments to report\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s %8s %9s %8s\n", "ENV", "TOTAL", "SENSITIVE", "REDACTED"))
	sb.WriteString(strings.Repeat("-", 50) + "\n")
	for _, e := range r.Entries {
		sb.WriteString(fmt.Sprintf("%-20s %8d %9d %8d\n",
			e.EnvName, e.TotalKeys, e.SensitiveKeys, e.RedactedKeys))
	}
	return sb.String()
}
