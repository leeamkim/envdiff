package diff

import (
	"fmt"
	"io"
	"sort"
)

// PrintMultiReport writes a summary of all pairwise diffs.
func PrintMultiReport(res MultiCompareResult, w io.Writer) {
	pairs := make([]string, 0, len(res.Pairs))
	for k := range res.Pairs {
		pairs = append(pairs, k)
	}
	sort.Strings(pairs)

	for _, key := range pairs {
		r := res.Pairs[key]
		s := Summarize(r)
		fmt.Fprintf(w, "[%s] missing_a=%d missing_b=%d mismatched=%d\n",
			key, s.MissingInA, s.MissingInB, s.Mismatched)
	}

	if !res.HasAnyDiff() {
		fmt.Fprintln(w, "all pairs match")
	}
}
