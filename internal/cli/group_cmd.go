package cli

import (
	"fmt"
	"io"
	"os"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunGroup parses two env files and prints diff entries grouped by key prefix.
func RunGroup(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff group <file1> <file2>")
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
	entries := diff.Flatten(result)
	groups := diff.GroupByPrefix(entries)

	if len(groups) == 0 {
		fmt.Fprintln(out, "no entries found")
		return nil
	}

	for _, g := range groups {
		fmt.Fprintf(out, "[%s]\n", g.Prefix)
		for _, e := range g.Entries {
			fmt.Fprintf(out, "  %-30s %s\n", e.Key, e.Status)
		}
	}

	return nil
}

// RunGroupMain is the entry point called from main for the group subcommand.
func RunGroupMain(args []string) {
	if err := RunGroup(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
