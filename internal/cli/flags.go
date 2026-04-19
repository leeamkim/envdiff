package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

// Usage prints CLI usage to w.
func Usage(w io.Writer) {
	fmt.Fprintln(w, `envdiff — compare .env files across environments

Usage:
  envdiff <file1> <file2> [flags]
  envdiff lint <file> [flags]
  envdiff merge <file1> <file2> [flags]
  envdiff redact <file> [flags]
  envdiff schema <file> [flags]
  envdiff template <template> <file> [flags]
  envdiff annotate <file1> <file2> [flags]
  envdiff group <file1> <file2> [flags]
  envdiff baseline <save|check> <file> [flags]
  envdiff watch <file1> <file2> [flags]

Flags:`)
	flag.CommandLine.SetOutput(w)
	flag.PrintDefaults()
}

// ParseFlags parses top-level CLI flags and returns remaining args.
type Flags struct {
	Strict   bool
	Format   string
	Ignore   string
	Interval time.Duration
}

func ParseFlags(args []string) (Flags, []string, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var f Flags
	fs.BoolVar(&f.Strict, "strict", false, "exit with non-zero status if any diff found")
	fs.StringVar(&f.Format, "format", "text", "output format: text, json, csv")
	fs.StringVar(&f.Ignore, "ignore", "", "path to ignore file")
	fs.DurationVar(&f.Interval, "interval", 2*time.Second, "polling interval for watch command")

	if err := fs.Parse(args); err != nil {
		return Flags{}, nil, err
	}
	return f, fs.Args(), nil
}
