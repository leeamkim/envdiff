package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunMulti compares all provided env files pairwise.
func RunMulti(args []string, out io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("multi requires at least 2 env files")
	}

	envs := make(map[string]map[string]string, len(args))
	for _, path := range args {
		m, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		envs[path] = m
	}

	res := diff.CompareMulti(envs)

	pairs := make([]string, 0, len(res.Pairs))
	for k := range res.Pairs {
		pairs = append(pairs, k)
	}
	sort.Strings(pairs)

	for _, key := range pairs {
		r := res.Pairs[key]
		fmt.Fprintf(out, "=== %s ===\n", key)
		diff.PrintReport(r, out)
	}

	if res.HasAnyDiff() {
		os.Exit(1)
	}
	return nil
}
