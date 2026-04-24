package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunInterpolate loads a .env file, expands ${VAR} references within it,
// prints the resolved values, and reports any unresolved references.
//
// Usage: envdiff interpolate <file>
func RunInterpolate(args []string, out io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff interpolate <file>")
	}

	filePath := args[0]
	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %q: %w", filePath, err)
	}

	resolved, issues := diff.Interpolate(env)

	// Print resolved key=value pairs in sorted order.
	keys := sortedStringKeys(resolved)
	for _, k := range keys {
		fmt.Fprintf(out, "%s=%s\n", k, resolved[k])
	}

	if len(issues) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, "unresolved references:")
		fmt.Fprint(out, diff.FormatInterpolateIssues(issues))
		return fmt.Errorf("%d unresolved interpolation reference(s)", len(issues))
	}

	return nil
}

// sortedStringKeys returns sorted keys of a string map.
func sortedStringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// reuse existing sort helper from the diff package via simple insertion sort
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}

// RunInterpolateMain is the entry point called from main dispatch.
func RunInterpolateMain(args []string) {
	if err := RunInterpolate(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
