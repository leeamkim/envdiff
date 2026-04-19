package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunGraph compares multiple named env files as a graph.
// Usage: envdiff graph name1=file1 name2=file2 edge1:edge2 ...
func RunGraph(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: graph <name=file>... <from:to>...")
	}

	g := diff.NewEnvGraph()

	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			name, file := parts[0], parts[1]
			vars, err := parser.ParseFile(file)
			if err != nil {
				return fmt.Errorf("parse %s: %w", file, err)
			}
			g.AddNode(name, vars)
		} else if strings.Contains(arg, ":") {
			parts := strings.SplitN(arg, ":", 2)
			g.AddEdge(parts[0], parts[1])
		} else {
			return fmt.Errorf("unrecognized argument: %s", arg)
		}
	}

	entries := diff.DiffGraph(g)
	if len(entries) == 0 {
		fmt.Fprintln(out, "no edges to compare")
		return nil
	}

	hasAny := false
	for _, e := range entries {
		if e.Result.HasDiff() {
			hasAny = true
		}
		fmt.Fprintln(out, e.String())
	}

	if hasAny {
		os.Exit(1)
	}
	return nil
}
