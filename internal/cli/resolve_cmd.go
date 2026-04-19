package cli

import (
	"fmt"
	"os"
	"sort"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunResolve merges two env files using a chosen strategy and prints the result.
func RunResolve(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: resolve <strategy> <fileA> <fileB>\nstrategies: prefer-a, prefer-b, require-match")
	}

	strategyStr := args[0]
	fileA := args[1]
	fileB := args[2]

	var strategy diff.ResolveStrategy
	switch strategyStr {
	case "prefer-a":
		strategy = diff.StrategyPreferA
	case "prefer-b":
		strategy = diff.StrategyPreferB
	case "require-match":
		strategy = diff.StrategyRequireMatch
	default:
		return fmt.Errorf("unknown strategy %q: use prefer-a, prefer-b, or require-match", strategyStr)
	}

	a, err := parser.ParseFile(fileA)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", fileA, err)
	}
	b, err := parser.ParseFile(fileB)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", fileB, err)
	}

	result := diff.Resolve(a, b, strategy)

	if len(result.Conflicts) > 0 {
		fmt.Fprintln(os.Stderr, "conflicts:")
		for _, c := range result.Conflicts {
			fmt.Fprintln(os.Stderr, " ", c)
		}
		return fmt.Errorf("%d conflict(s) found", len(result.Conflicts))
	}

	keys := make([]string, 0, len(result.Resolved))
	for k := range result.Resolved {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s=%s\n", k, result.Resolved[k])
	}
	return nil
}
