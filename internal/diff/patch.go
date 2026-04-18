package diff

import (
	"fmt"
	"sort"
	"strings"
)

// PatchEntry represents a single change instruction.
type PatchEntry struct {
	Key    string
	Action string // "add", "remove", "update"
	Value  string
}

// GeneratePatch produces a list of patch entries to transform env map a into b.
func GeneratePatch(a, b map[string]string) []PatchEntry {
	patch := []PatchEntry{}

	for k, vb := range b {
		va, exists := a[k]
		if !exists {
			patch = append(patch, PatchEntry{Key: k, Action: "add", Value: vb})
		} else if va != vb {
			patch = append(patch, PatchEntry{Key: k, Action: "update", Value: vb})
		}
	}

	for k := range a {
		if _, exists := b[k]; !exists {
			patch = append(patch, PatchEntry{Key: k, Action: "remove"})
		}
	}

	sort.Slice(patch, func(i, j int) bool {
		return patch[i].Key < patch[j].Key
	})
	return patch
}

// FormatPatch returns a human-readable patch string.
func FormatPatch(patch []PatchEntry) string {
	var sb strings.Builder
	for _, e := range patch {
		switch e.Action {
		case "add":
			fmt.Fprintf(&sb, "+ %s=%s\n", e.Key, e.Value)
		case "remove":
			fmt.Fprintf(&sb, "- %s\n", e.Key)
		case "update":
			fmt.Fprintf(&sb, "~ %s=%s\n", e.Key, e.Value)
		}
	}
	return sb.String()
}
