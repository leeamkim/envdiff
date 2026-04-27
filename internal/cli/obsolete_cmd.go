package cli

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunObsolete compares an env file against a reference and reports obsolete keys.
// Usage: envdiff obsolete <env-file> <reference-file>
func RunObsolete(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff obsolete <env-file> <reference-file>")
	}

	envFile := args[0]
	refFile := args[1]

	env, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to parse env file %q: %w", envFile, err)
	}

	ref, err := parser.ParseFile(refFile)
	if err != nil {
		return fmt.Errorf("failed to parse reference file %q: %w", refFile, err)
	}

	issues := diff.CheckObsolete(env, ref)

	if len(issues) == 0 {
		fmt.Fprintln(os.Stdout, "no obsolete keys found")
		return nil
	}

	fmt.Fprintf(os.Stdout, "obsolete keys in %s (not present in %s):\n", envFile, refFile)
	for _, issue := range issues {
		fmt.Fprintf(os.Stdout, "  %s\n", issue.String())
	}
	return nil
}
