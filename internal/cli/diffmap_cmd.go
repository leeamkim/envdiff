package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunDiffMap compares N env files and prints a key-by-env consistency table.
// Usage: envdiff diffmap <name=file> <name=file> ...
func RunDiffMap(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("diffmap requires at least 2 name=file arguments")
	}

	envs := make(map[string]map[string]string, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid argument %q: expected name=file", arg)
		}
		name, path := parts[0], parts[1]
		vars, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}
		envs[name] = vars
	}

	m := diff.BuildEnvDiffMap(envs)
	printDiffMap(m, out)
	return nil
}

func printDiffMap(m *diff.EnvDiffMap, out io.Writer) {
	// header
	fmt.Fprintf(out, "%-30s", "KEY")
	for _, env := range m.Envs {
		fmt.Fprintf(out, "  %-14s", env)
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, strings.Repeat("-", 30+16*len(m.Envs)))

	for _, key := range m.Keys {
		marker := " "
		if !m.Consistent(key) {
			marker = "!"
		}
		fmt.Fprintf(out, "%s %-29s", marker, key)
		for _, env := range m.Envs {
			v := m.Cells[key][env]
			if v == "" {
				v = "<missing>"
			}
			if len(v) > 13 {
				v = v[:10] + "..."
			}
			fmt.Fprintf(out, "  %-14s", v)
		}
		fmt.Fprintln(out)
	}
}

func RunDiffMapMain(args []string) {
	if err := RunDiffMap(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
