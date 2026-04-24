package diff

import (
	"fmt"
	"strings"
)

// FormatChainCSV returns a CSV representation of the chain result.
// Columns: key, env1, env2, ...
func FormatChainCSV(r ChainResult) string {
	if len(r.Entries) == 0 {
		return ""
	}
	var sb strings.Builder
	// header
	sb.WriteString("key")
	for _, env := range r.Chain {
		sb.WriteString(",")
		sb.WriteString(env)
	}
	sb.WriteString("\n")
	// rows
	for _, e := range r.Entries {
		sb.WriteString(e.Key)
		for _, cv := range e.Values {
			sb.WriteString(",")
			if cv.Present {
				sb.WriteString(cv.Value)
			} else {
				sb.WriteString("(missing)")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// SummarizeChain returns a short summary: how many keys have at least one override.
func SummarizeChain(r ChainResult) string {
	overrides := 0
	missing := 0
	for _, e := range r.Entries {
		hasOverride := false
		for _, cv := range e.Values {
			if cv.Override {
				hasOverride = true
			}
			if !cv.Present {
				missing++
			}
		}
		if hasOverride {
			overrides++
		}
	}
	return fmt.Sprintf("keys: %d  overrides: %d  missing-slots: %d",
		len(r.Entries), overrides, missing)
}
