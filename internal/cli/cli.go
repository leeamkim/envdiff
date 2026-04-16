package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// Run is the entry point for the CLI.
func Run(args []string) error {
	files, opts, err := parseArgs(args)
	if err != nil {
		return err
	}
	return run(files[0], files[1], opts)
}

type options struct {
	strict      bool
	onlyMissing bool
	onlyMismatch bool
}

func parseArgs(args []string) ([2]string, options, error) {
	var files [2]string
	var opts options
	positional := []string{}
	for _, a := range args {
		switch a {
		case "--strict":
			opts.strict = true
		case "--only-missing":
			opts.onlyMissing = true
		case "--only-mismatch":
			opts.onlyMismatch = true
		default:
			positional = append(positional, a)
		}
	}
	if len(positional) < 2 {
		return files, opts, errors.New("usage: envdiff <file1> <file2> [--strict] [--only-missing] [--only-mismatch]")
	}
	files[0], files[1] = positional[0], positional[1]
	return files, opts, nil
}

func run(fileA, fileB string, opts options) error {
	envA, err := parser.ParseFile(fileA)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", fileA, err)
	}
	envB, err := parser.ParseFile(fileB)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", fileB, err)
	}

	results := diff.Compare(envA, envB)
	results = diff.Filter(results, diff.FilterOptions{
		OnlyMissing:   opts.onlyMissing,
		OnlyMismatched: opts.onlyMismatch,
	})

	diff.PrintReport(os.Stdout, results, fileA, fileB)

	summary := diff.Summarize(results)
	fmt.Fprintln(os.Stdout, summary.String())

	if opts.strict && summary.HasDiff() {
		return errors.New("differences found (strict mode)")
	}
	return nil
}
