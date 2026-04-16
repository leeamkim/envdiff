package diff

import (
	"fmt"
	"io"
	"sort"
)

// PrintReport writes a human-readable diff report to w.
func PrintReport(w io.Writer, r Result, leftName, rightName string) {
	if !r.HasDiff() {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	if len(r.MissingInRight) > 0 {
		sort.Strings(r.MissingInRight)
		fmt.Fprintf(w, "Keys missing in %s:\n", rightName)
		for _, k := range r.MissingInRight {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(r.MissingInLeft) > 0 {
		sort.Strings(r.MissingInLeft)
		fmt.Fprintf(w, "Keys missing in %s:\n", leftName)
		for _, k := range r.MissingInLeft {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(r.Mismatched) > 0 {
		keys := make([]string, 0, len(r.Mismatched))
		for k := range r.Mismatched {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Fprintln(w, "Mismatched values:")
		for _, k := range keys {
			v := r.Mismatched[k]
			fmt.Fprintf(w, "  ~ %s: %q (%s) vs %q (%s)\n", k, v[0], leftName, v[1], rightName)
		}
	}
}
