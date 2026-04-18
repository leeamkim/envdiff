package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunMerge executes the merge subcommand writing result to w.
func RunMerge(args []string, w io.Writer) int {
	if len(args) < 2 {
		fmt.Fprintln(w, "usage: envdiff merge <fileA> <fileB> [--prefer-a]")
		return 1
	}

	fileA, fileB := args[0], args[1]
	preferA := len(args) >= 3 && args[2] == "--prefer-a"

	a, err := parser.ParseFile(fileA)
	if err != nil {
		fmt.Fprintf(w, "error reading %s: %v\n", fileA, err)
		return 1
	}

	b, err := parser.ParseFile(fileB)
	if err != nil {
		fmt.Fprintf(w, "error reading %s: %v\n", fileB, err)
		return 1
	}

	result := diff.Merge(a, b, preferA)

	if len(result.Conflicts) > 0 {
		fmt.Fprintf(w, "# conflicts resolved (%d):\n", len(result.Conflicts))
		for _, k := range result.Conflicts {
			fmt.Fprintf(w, "#   %s\n", k)
		}
	}

	keys := make([]string, 0, len(result.Merged))
	for k := range result.Merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "%s=%s\n", k, result.Merged[k])
	}

	_ = os.Stdout
	return 0
}
