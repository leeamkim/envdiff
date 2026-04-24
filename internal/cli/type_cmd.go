package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunType runs the type-inference and consistency check command.
// Usage: envdiff type <name=file> [<name=file> ...]
func RunType(args []string, out io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff type <name=file> [<name=file> ...]")
	}

	envs := map[string]map[string]string{}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid argument %q: expected name=file", arg)
		}
		name, path := parts[0], parts[1]
		vars, err := parser.ParseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %q: %w", path, err)
		}
		envs[name] = vars
	}

	// Print per-env type annotations
	for _, name := range sortedEnvNames(envs) {
		entries := diff.InferTypes(envs[name])
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
		fmt.Fprintf(out, "[%s]\n", name)
		for _, e := range entries {
			fmt.Fprintf(out, "  %-30s %s\n", e.Key, e.Type)
		}
	}

	// Check consistency across envs if more than one provided
	if len(envs) > 1 {
		issues := diff.CheckTypeConsistency(envs)
		if len(issues) == 0 {
			fmt.Fprintln(out, "\nAll shared keys have consistent types.")
		} else {
			fmt.Fprintf(out, "\n%d type inconsistency(ies) found:\n", len(issues))
			sort.Slice(issues, func(i, j int) bool {
				return issues[i].Key < issues[j].Key
			})
			for _, issue := range issues {
				fmt.Fprintf(out, "  %s\n", issue.String())
			}
			os.Exit(1)
		}
	}
	return nil
}

func sortedEnvNames(envs map[string]map[string]string) []string {
	names := make([]string, 0, len(envs))
	for n := range envs {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
