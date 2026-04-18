package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// Run is the entry point for the CLI.
func Run(args []string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}
	return run(opts)
}

type options struct {
	fileA      string
	fileB      string
	strict     bool
	ignoreFile string
}

func parseArgs(args []string) (options, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	strict := fs.Bool("strict", false, "exit with non-zero status if any diff found")
	ignoreFile := fs.String("ignore", "", "path to file listing keys to ignore")
	if err := fs.Parse(args); err != nil {
		return options{}, err
	}
	if fs.NArg() < 2 {
		return options{}, errors.New("usage: envdiff [--strict] [--ignore FILE] <fileA> <fileB>")
	}
	return options{
		fileA:      fs.Arg(0),
		fileB:      fs.Arg(1),
		strict:     *strict,
		ignoreFile: *ignoreFile,
	}, nil
}

func run(opts options) error {
	envA, err := parser.ParseFile(opts.fileA)
	if err != nil {
		return fmt.Errorf("reading %s: %w", opts.fileA, err)
	}
	envB, err := parser.ParseFile(opts.fileB)
	if err != nil {
		return fmt.Errorf("reading %s: %w", opts.fileB, err)
	}

	result := diff.Compare(envA, envB)

	if opts.ignoreFile != "" {
		keys, err := parser.ParseIgnoreFile(opts.ignoreFile)
		if err != nil {
			return fmt.Errorf("reading ignore file: %w", err)
		}
		il := diff.NewIgnoreList(keys)
		result = il.Apply(result)
	}

	diff.PrintReport(os.Stdout, opts.fileA, opts.fileB, result)

	if opts.strict && diff.Summarize(result).HasDiff() {
		return errors.New("differences found")
	}
	return nil
}
