package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunFlatten loads multiple .env files and prints a flat list of all key=value entries with their source env.
func RunFlatten(args []string, out io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff flatten <file1> [file2 ...]")
	}

	envs := make(map[string]map[string]string)
	for _, path := range args {
		vars, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}
		envs[path] = vars
	}

	entries := diff.FlattenEnvs(envs)
	grouped := diff.GroupFlatByKey(entries)

	keys := make([]string, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		for _, e := range grouped[key] {
			fmt.Fprintf(out, "[%s] %s=%s\n", e.Env, e.Key, e.Value)
		}
	}
	return nil
}

// RunFlattenMain is the entry point for the flatten subcommand.
func RunFlattenMain(args []string) {
	if err := RunFlatten(args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
