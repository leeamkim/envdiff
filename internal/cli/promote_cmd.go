package cli

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunPromote promotes keys from one .env file into another.
// Usage: envdiff promote <src> <dest> [--overwrite] [--keys KEY1,KEY2]
func RunPromote(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff promote <src> <dest> [--overwrite]")
	}

	srcPath := args[0]
	destPath := args[1]

	overwrite := false
	var keys []string

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "--overwrite":
			overwrite = true
		case "--keys":
			if i+1 < len(args) {
				i++
				keys = splitCSV(args[i])
			}
		}
	}

	src, err := parser.ParseFile(srcPath)
	if err != nil {
		return fmt.Errorf("parse src: %w", err)
	}

	dest, err := parser.ParseFile(destPath)
	if err != nil {
		return fmt.Errorf("parse dest: %w", err)
	}

	var opts []diff.PromoteOption
	if overwrite {
		opts = append(opts, diff.PromoteOverwrite())
	}
	if len(keys) > 0 {
		opts = append(opts, diff.PromoteKeys(keys...))
	}

	result := diff.Promote(src, dest, opts...)

	if len(result.Conflicts) > 0 {
		fmt.Fprintln(os.Stderr, "conflicts (use --overwrite to force):")
		for k := range result.Conflicts {
			fmt.Fprintf(os.Stderr, "  %s\n", k)
		}
	}

	fmt.Println(diff.FormatPromoteResult(result))
	return nil
}

func splitCSV(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if tok := s[start:i]; tok != "" {
				out = append(out, tok)
			}
			start = i + 1
		}
	}
	return out
}
