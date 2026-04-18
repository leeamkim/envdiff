package cli

import (
	"flag"
	"fmt"
	"io"
)

// Flags holds all parsed command-line flags for envdiff.
type Flags struct {
	FileA      string
	FileB      string
	Strict     bool
	OnlyMissing    bool
	OnlyMismatched bool
	IgnoreFile string
	Format     string
	Output     string
	ShowStats  bool
	ShowSummary bool
}

// Usage prints a formatted help message to the given writer.
func Usage(w io.Writer) {
	fmt.Fprintf(w, `envdiff — compare .env files across environments

Usage:
  envdiff [flags] <file-a> <file-b>

Flags:
  --strict            Exit with non-zero status if any diff is found
  --only-missing      Show only missing keys
  --only-mismatched   Show only mismatched values
  --ignore-file       Path to file with keys to ignore
  --format            Output format: text (default), json, csv
  --output            Write output to file instead of stdout
  --stats             Show diff statistics
  --summary           Show a short summary line
  --help              Show this help message

Examples:
  envdiff .env.development .env.production
  envdiff --strict --format json .env .env.production
  envdiff --ignore-file .envdiffignore .env .env.staging
`)
}

// ParseFlags parses command-line arguments into a Flags struct.
// It returns the remaining positional arguments alongside the parsed flags.
func ParseFlags(args []string) (*Flags, []string, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)

	f := &Flags{}

	fs.BoolVar(&f.Strict, "strict", false, "exit non-zero if any diff found")
	fs.BoolVar(&f.OnlyMissing, "only-missing", false, "show only missing keys")
	fs.BoolVar(&f.OnlyMismatched, "only-mismatched", false, "show only mismatched values")
	fs.StringVar(&f.IgnoreFile, "ignore-file", "", "path to ignore file")
	fs.StringVar(&f.Format, "format", "text", "output format: text, json, csv")
	fs.StringVar(&f.Output, "output", "", "write output to file")
	fs.BoolVar(&f.ShowStats, "stats", false, "show diff statistics")
	fs.BoolVar(&f.ShowSummary, "summary", false, "show summary line")

	if err := fs.Parse(args); err != nil {
		return nil, nil, err
	}

	return f, fs.Args(), nil
}
