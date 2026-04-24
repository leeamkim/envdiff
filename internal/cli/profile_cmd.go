package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunProfile compares multiple named env files and prints a cross-profile table.
// Usage: envdiff profile name1=file1.env name2=file2.env ...
func RunProfile(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("profile requires at least 2 arguments: name=file pairs")
	}

	envs := make(map[string]map[string]string, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid profile argument %q: expected name=file", arg)
		}
		name, file := parts[0], parts[1]
		env, err := parser.ParseFile(file)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", file, err)
		}
		envs[name] = env
	}

	report := diff.BuildProfileReport(envs)
	fmt.Fprint(out, diff.FormatProfileReport(report))
	return nil
}

// RunProfileMain is the entry point called from main dispatch.
func RunProfileMain(args []string) {
	if err := RunProfile(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
