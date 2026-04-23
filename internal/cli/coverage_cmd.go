package cli

import (
	"fmt"
	"os"

	"github.com/your-org/envdiff/internal/diff"
	"github.com/your-org/envdiff/internal/parser"
)

// RunCoverage compares one or more target env files against a reference file
// and prints a coverage report.
//
// Usage: envdiff coverage <reference> <target> [<target2> ...]
func RunCoverage(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff coverage <reference> <target> [<target2> ...]")
	}

	refPath := args[0]
	targetPaths := args[1:]

	ref, err := parser.ParseFile(refPath)
	if err != nil {
		return fmt.Errorf("reading reference %s: %w", refPath, err)
	}

	results := make(map[string]diff.CoverageResult, len(targetPaths))
	for _, tp := range targetPaths {
		tgt, err := parser.ParseFile(tp)
		if err != nil {
			return fmt.Errorf("reading target %s: %w", tp, err)
		}
		results[tp] = diff.ComputeCoverage(ref, tgt)
	}

	fmt.Fprint(os.Stdout, diff.FormatMultiCoverageReport(results))

	// Exit non-zero if any environment has grade F.
	for _, r := range results {
		if r.Grade == "F" {
			return fmt.Errorf("one or more environments have failing coverage")
		}
	}
	return nil
}
