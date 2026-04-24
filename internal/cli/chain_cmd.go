package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/your/envdiff/internal/diff"
	"github.com/your/envdiff/internal/parser"
)

// RunChain runs the chain command: envdiff chain env1=file1 env2=file2 ...
// It shows how each key's value propagates (or is overridden) across an ordered
// list of environments.
func RunChain(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("chain requires at least 2 env=file arguments")
	}

	chain := make([]string, 0, len(args))
	envs := make(map[string]map[string]string, len(args))

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid argument %q: expected name=file", arg)
		}
		name, file := parts[0], parts[1]
		vars, err := parser.ParseFile(file)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", file, err)
		}
		chain = append(chain, name)
		envs[name] = vars
	}

	result := diff.BuildChain(chain, envs)
	report := diff.FormatChainReport(result)
	fmt.Fprint(out, report)
	return nil
}

// RunChainMain is the entry point used by main dispatch.
func RunChainMain(args []string) {
	if err := RunChain(args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
