package cli

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunFormat formats a .env file and prints it to stdout using the given style.
func RunFormat(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff format <file> [--style=plain|export|quoted] [--prefix=PREFIX] [--sort]")
	}

	filePath := args[0]
	style := diff.FormatStylePlain
	prefix := ""
	sortKeys := false

	for _, arg := range args[1:] {
		switch {
		case arg == "--sort":
			sortKeys = true
		case len(arg) > 8 && arg[:8] == "--style=":
			switch diff.FormatStyle(arg[8:]) {
			case diff.FormatStyleExport, diff.FormatStylePlain, diff.FormatStyleQuoted:
				style = diff.FormatStyle(arg[8:])
			default:
				return fmt.Errorf("unknown style %q: use plain, export, or quoted", arg[8:])
			}
		case len(arg) > 9 && arg[:9] == "--prefix=":
			prefix = arg[9:]
		}
	}

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	out := diff.FormatEnv(env, diff.FormatOptions{
		Style:    style,
		SortKeys: sortKeys,
		Prefix:   prefix,
	})

	fmt.Fprint(os.Stdout, out)
	return nil
}
