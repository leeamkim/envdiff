package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// Config holds parsed CLI options.
type Config struct {
	FileA  string
	FileB  string
	Strict bool
	Output io.Writer
}

// Run parses arguments and executes the comparison.
func Run(args []string) error {
	cfg, err := parseArgs(args)
	if err != nil {
		return err
	}
	return run(cfg)
}

func parseArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	strict := fs.Bool("strict", false, "exit with non-zero status if differences found")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() != 2 {
		return nil, errors.New("usage: envdiff [--strict] <file-a> <file-b>")
	}

	return &Config{
		FileA:  fs.Arg(0),
		FileB:  fs.Arg(1),
		Strict: *strict,
		Output: os.Stdout,
	}, nil
}

func run(cfg *Config) error {
	envA, err := parser.ParseFile(cfg.FileA)
	if err != nil {
		return fmt.Errorf("reading %s: %w", cfg.FileA, err)
	}

	envB, err := parser.ParseFile(cfg.FileB)
	if err != nil {
		return fmt.Errorf("reading %s: %w", cfg.FileB, err)
	}

	result := diff.Compare(envA, envB)
	diff.PrintReport(cfg.Output, result, cfg.FileA, cfg.FileB)

	if cfg.Strict && result.HasDiff() {
		return fmt.Errorf("differences found between %s and %s", cfg.FileA, cfg.FileB)
	}
	return nil
}
