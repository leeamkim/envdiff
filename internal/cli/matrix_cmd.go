package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunMatrix compares multiple .env files and prints a cross-environment matrix.
func RunMatrix(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff matrix <file1> <file2> [file3...]")
	}

	envs := map[string]map[string]string{}
	for _, path := range args {
		vars, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}
		envs[path] = vars
	}

	matrix := diff.BuildMatrix(envs)

	if len(matrix.Keys) == 0 {
		fmt.Fprintln(out, "No keys found.")
		return nil
	}

	base := matrix.EnvNames[0]
	others := matrix.EnvNames[1:]

	// Header
	header := fmt.Sprintf("%-30s", "KEY")
	for _, name := range others {
		short := shortName(name)
		header += fmt.Sprintf("  %-12s", short)
	}
	fmt.Fprintln(out, header)
	fmt.Fprintln(out, strings.Repeat("-", len(header)))

	for _, key := range matrix.Keys {
		row := fmt.Sprintf("%-30s", key)
		for _, name := range others {
			pair := base + ":" + name
			cell := matrix.Status[key][pair]
			var symbol string
			switch cell.Status {
			case "match":
				symbol = "OK"
			case "mismatch":
				symbol = "MISMATCH"
			case "missing":
				symbol = "MISSING"
			default:
				symbol = "?"
			}
			row += fmt.Sprintf("  %-12s", symbol)
		}
		fmt.Fprintln(out, row)
	}
	return nil
}

func shortName(path string) string {
	base := path
	if idx := strings.LastIndexAny(path, "/\\"); idx >= 0 {
		base = path[idx+1:]
	}
	return base
}

func RunMatrixMain() {
	if err := RunMatrix(os.Args[2:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
