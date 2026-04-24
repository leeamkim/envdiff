package diff

import (
	"fmt"
	"sort"
	"strings"
)

// MaskEntry represents a single masked key result.
type MaskEntry struct {
	Key      string
	Original string
	Masked   string
}

// MaskOptions controls how values are masked.
type MaskOptions struct {
	// ShowLength reveals the value length in the mask placeholder.
	ShowLength bool
	// VisibleChars shows the first N characters before masking.
	VisibleChars int
}

// MaskEnv returns a copy of env with sensitive values masked.
// Keys matching isSensitive are redacted according to opts.
func MaskEnv(env map[string]string, opts MaskOptions) []MaskEntry {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]MaskEntry, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		var masked string
		if isSensitive(k) {
			masked = buildMask(v, opts)
		} else {
			masked = v
		}
		entries = append(entries, MaskEntry{Key: k, Original: v, Masked: masked})
	}
	return entries
}

func buildMask(value string, opts MaskOptions) string {
	if len(value) == 0 {
		return ""
	}
	visible := ""
	if opts.VisibleChars > 0 && opts.VisibleChars < len(value) {
		visible = value[:opts.VisibleChars]
	}
	if opts.ShowLength {
		return fmt.Sprintf("%s****(%d)", visible, len(value))
	}
	return visible + strings.Repeat("*", 4)
}

// FormatMaskEntries returns a human-readable table of masked entries.
func FormatMaskEntries(entries []MaskEntry) string {
	if len(entries) == 0 {
		return "no entries\n"
	}
	var sb strings.Builder
	for _, e := range entries {
		if e.Original != e.Masked {
			fmt.Fprintf(&sb, "%s = %s [masked]\n", e.Key, e.Masked)
		} else {
			fmt.Fprintf(&sb, "%s = %s\n", e.Key, e.Masked)
		}
	}
	return sb.String()
}
