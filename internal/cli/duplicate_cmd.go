package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunDuplicate compares multiple .env files for keys that exist in all files
// but carry different values, reporting each conflict.
//
// Usage: envdiff duplicate <name=file> [<name=file> ...]
func RunDuplicate(args []string, w io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("duplicate requires at least two name=file arguments")
	}

	envs := make(map[string]map[string]string, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid argument %q: expected name=file", arg)
		}
		name, file := parts[0], parts[1]
		m, err := parser.ParseFile(file)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", file, err)
		}
		envs[name] = m
	}

	issues := diff.FindDuplicateConflicts(envs)
	if len(issues) == 0 {
		fmt.Fprintln(w, "no duplicate conflicts found")
		return nil
	}

	fmt.Fprintln(w, diff.FormatDuplicateIssues(issues))
	return nil
}

// RunDuplicateMain is the entry point called from main dispatch.
func RunDuplicateMain(args []string) {
	if err := RunDuplicate(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
