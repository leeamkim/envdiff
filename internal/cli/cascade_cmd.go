package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunCascade compares a key across multiple env files.
// Usage: envdiff cascade file1=dev file2=staging file3=prod
func RunCascade(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("cascade requires at least 2 env files (name=path)")
	}

	var names []string
	envs := map[string]map[string]string{}

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid argument %q: expected name=path", arg)
		}
		name, path := parts[0], parts[1]
		m, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}
		names = append(names, name)
		envs[name] = m
	}

	result := diff.Cascade(names, envs)

	sort.Slice(result.Entries, func(i, j int) bool {
		return result.Entries[i].Key < result.Entries[j].Key
	})

	driftCount := 0
	for _, e := range result.Entries {
		status := "ok"
		if e.Drift {
			status = "DRIFT"
			driftCount++
		}
		fmt.Fprintf(os.Stdout, "%-30s %s\n", e.Key, status)
	}

	fmt.Fprintf(os.Stdout, "\n%d key(s) with drift across [%s]\n", driftCount, strings.Join(names, ", "))
	return nil
}
