package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/your-org/envdiff/internal/diff"
	"github.com/your-org/envdiff/internal/parser"
)

// RunMask parses an env file and prints its contents with sensitive values masked.
// Usage: envdiff mask <file> [--show-length] [--visible=N]
func RunMask(args []string, out io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff mask <file> [--show-length] [--visible=N]")
	}

	filePath := args[0]
	opts := diff.MaskOptions{}

	for _, arg := range args[1:] {
		switch {
		case arg == "--show-length":
			opts.ShowLength = true
		case strings.HasPrefix(arg, "--visible="):
			nStr := strings.TrimPrefix(arg, "--visible=")
			n, err := strconv.Atoi(nStr)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid --visible value: %q", nStr)
			}
			opts.VisibleChars = n
		default:
			return fmt.Errorf("unknown flag: %s", arg)
		}
	}

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	entries := diff.MaskEnv(env, opts)
	fmt.Fprint(out, diff.FormatMaskEntries(entries))
	return nil
}

// RunMaskMain is the entry point called from main dispatch.
func RunMaskMain(args []string) {
	if err := RunMask(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
