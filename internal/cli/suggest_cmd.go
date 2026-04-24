package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunSuggest compares keys in the source env against a reference env and
// prints rename suggestions for keys that are absent in the reference but
// have a likely match via case or separator normalization.
//
// Usage: envdiff suggest <source.env> <reference.env>
func RunSuggest(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff suggest <source.env> <reference.env>")
	}

	src, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("parsing source: %w", err)
	}

	ref, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("parsing reference: %w", err)
	}

	issues := diff.SuggestRenames(src, ref)
	if len(issues) == 0 {
		fmt.Fprintln(out, "no rename suggestions")
		return nil
	}

	fmt.Fprintf(out, "%d rename suggestion(s):\n", len(issues))
	for _, iss := range issues {
		fmt.Fprintf(out, "  %s\n", iss.String())
	}
	return nil
}

// RunSuggestMain is the entry point called from main dispatch.
func RunSuggestMain(args []string) {
	if err := RunSuggest(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
