package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ProfileEntry represents a single key's presence across named profiles.
type ProfileEntry struct {
	Key      string
	Profiles map[string]string // profile name -> value
}

// ProfileReport holds the full cross-profile analysis.
type ProfileReport struct {
	Profiles []string
	Entries  []ProfileEntry
}

// BuildProfileReport compares a set of named env maps and returns a ProfileReport
// showing each key's value (or absence) across all profiles.
func BuildProfileReport(envs map[string]map[string]string) ProfileReport {
	profileNames := make([]string, 0, len(envs))
	for name := range envs {
		profileNames = append(profileNames, name)
	}
	sort.Strings(profileNames)

	keySet := map[string]struct{}{}
	for _, env := range envs {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]ProfileEntry, 0, len(keys))
	for _, key := range keys {
		profiles := make(map[string]string, len(profileNames))
		for _, p := range profileNames {
			if v, ok := envs[p][key]; ok {
				profiles[p] = v
			} else {
				profiles[p] = ""
			}
		}
		entries = append(entries, ProfileEntry{Key: key, Profiles: profiles})
	}

	return ProfileReport{Profiles: profileNames, Entries: entries}
}

// FormatProfileReport returns a human-readable table of the profile report.
func FormatProfileReport(r ProfileReport) string {
	if len(r.Entries) == 0 {
		return "no keys found across profiles\n"
	}
	var sb strings.Builder
	header := fmt.Sprintf("%-30s", "KEY")
	for _, p := range r.Profiles {
		header += fmt.Sprintf("  %-20s", p)
	}
	sb.WriteString(header + "\n")
	sb.WriteString(strings.Repeat("-", len(header)) + "\n")
	for _, e := range r.Entries {
		line := fmt.Sprintf("%-30s", e.Key)
		for _, p := range r.Profiles {
			v := e.Profiles[p]
			if v == "" {
				v = "(missing)"
			}
			line += fmt.Sprintf("  %-20s", v)
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}
