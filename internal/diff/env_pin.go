package diff

import (
	"fmt"
	"sort"
)

// PinEntry represents a key whose value differs from the pinned (expected) value.
type PinEntry struct {
	Key      string
	Pinned   string
	Actual   string
	Missing  bool
}

func (e PinEntry) String() string {
	if e.Missing {
		return fmt.Sprintf("%s: pinned=%q (missing in env)", e.Key, e.Pinned)
	}
	return fmt.Sprintf("%s: pinned=%q actual=%q", e.Key, e.Pinned, e.Actual)
}

// CheckPins compares an env map against a map of pinned key=value pairs.
// It returns entries where the actual value differs or the key is absent.
func CheckPins(env map[string]string, pins map[string]string) []PinEntry {
	var issues []PinEntry
	for key, pinned := range pins {
		actual, ok := env[key]
		if !ok {
			issues = append(issues, PinEntry{Key: key, Pinned: pinned, Missing: true})
		} else if actual != pinned {
			issues = append(issues, PinEntry{Key: key, Pinned: pinned, Actual: actual})
		}
	}
	sort.Slice(issues, func(i, j int) bool {
		return issues[i].Key < issues[j].Key
	})
	return issues
}

// FormatPinIssues returns a human-readable summary of pin violations.
func FormatPinIssues(issues []PinEntry) string {
	if len(issues) == 0 {
		return "all pinned values match"
	}
	out := fmt.Sprintf("%d pin violation(s):\n", len(issues))
	for _, e := range issues {
		out += "  " + e.String() + "\n"
	}
	return out
}
