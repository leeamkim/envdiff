package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Flags holds all parsed CLI flags.
type Flags struct {
	FileA       string
	FileB       string
	Strict      bool
	IgnoreFile  string
	Format      string
	OnlyMissing bool
	OnlyMismatch bool
	// Lint subcommand flags
	LintFiles     []string
	LintNoEmpty   bool
	LintNoPlaceholder bool
	LintUppercase bool
}

// Usage prints usage information to w.
func Usage(w io.Writer) {
	fmt.Fprintln(w, "Usage: envdiff [options] <file-a> <file-b>")
	fmt.Fprintln(w, "       envdiff lint [--no-empty] [--no-placeholder] [--uppercase] <files...>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  --strict          Exit with non-zero status if any diff found")
	fmt.Fprintln(w, "  --ignore FILE     Path to ignore file")
	fmt.Fprintln(w, "  --format FORMAT   Output format: text (default), json, csv")
	fmt.Fprintln(w, "  --only-missing    Show only missing keys")
	fmt.Fprintln(w, "  --only-mismatch   Show only mismatched keys")
}

// ParseFlags parses os.Args and returns a Flags struct.
func ParseFlags(args []string) (Flags, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var f Flags
	fs.BoolVar(&f.Strict, "strict", false, "exit non-zero if diff found")
	fs.StringVar(&f.IgnoreFile, "ignore", "", "path to ignore file")
	fs.StringVar(&f.Format, "format", "text", "output format: text, json, csv")
	fs.BoolVar(&f.OnlyMissing, "only-missing", false, "show only missing keys")
	fs.BoolVar(&f.OnlyMismatch, "only-mismatch", false, "show only mismatched keys")
	fs.BoolVar(&f.LintNoEmpty, "no-empty", false, "lint: flag empty values")
	fs.BoolVar(&f.LintNoPlaceholder, "no-placeholder", false, "lint: flag placeholder values")
	fs.BoolVar(&f.LintUppercase, "uppercase", false, "lint: flag lowercase keys")

	if err := fs.Parse(args); err != nil {
		return f, err
	}

	remaining := fs.Args()
	if len(remaining) >= 2 {
		f.FileA = remaining[0]
		f.FileB = remaining[1]
	} else if len(remaining) == 1 && remaining[0] == "lint" {
		// handled upstream
	} else {
		f.LintFiles = remaining
	}
	return f, nil
}
