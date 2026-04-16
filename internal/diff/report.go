package diff

import (
	"fmt"
	"io"
	"sort"
)

// PrintReport writes a human-readable diff report to w.
func PrintReport(w io.Writer, r Result, labelA, labelB string) {
	if !r.HasDiff() {
		fmt.Fprintln(w, "✓ No differences found.")
		return
	}

	if len(r.MissingInB) > 0 {
		keys := sortedKeys(r.MissingInB)
		fmt.Fprintf(w, "\nKeys in %s but missing in %s:\n", labelA, labelB)
		for _, k := range keys {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(r.MissingInA) > 0 {
		keys := sortedKeys(r.MissingInA)
		fmt.Fprintf(w, "\nKeys in %s but missing in %s:\n", labelB, labelA)
		for _, k := range keys {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(r.Mismatched) > 0 {
		fmt.Fprintf(w, "\nMismatched values:\n")
		keys := make([]string, 0, len(r.Mismatched))
		for k := range r.Mismatched {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			pair := r.Mismatched[k]
			fmt.Fprintf(w, "  ~ %s: %q (%s) vs %q (%s)\n", k, pair[0], labelA, pair[1], labelB)
		}
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
