package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunAnnotate compares two env files and prints per-key annotations.
func RunAnnotate(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff annotate <file1> <file2>")
	}

	a, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[0], err)
	}
	b, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("reading %s: %w", args[1], err)
	}

	result := diff.Compare(a, b)
	ar := diff.Annotate(result, []func(diff.DiffEntry) string{
		diff.AnnotatorMissingInB,
		diff.AnnotatorMissingInA,
		diff.AnnotatorMismatched,
	})

	if len(ar.Annotations) == 0 {
		fmt.Fprintln(out, "no issues found")
		return nil
	}

	keys := make([]string, 0, len(ar.Annotations))
	for k := range ar.Annotations {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(out, "%-30s %s\n", k, ar.Annotations[k])
	}
	return nil
}

// RunAnnotateMain is the entry point wired into main.
func RunAnnotateMain(args []string) {
	if err := RunAnnotate(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
