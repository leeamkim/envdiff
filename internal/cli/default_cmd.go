package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunDefault checks one or more .env files for keys with default/placeholder values.
// Usage: envdiff default [--pins KEY=VAL,...] <file> [file2 ...]
func RunDefault(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff default [--pins KEY=VAL,...] <file> [file2 ...]")
	}

	var pinPairs []string
	files := args

	if len(args) >= 2 && args[0] == "--pins" {
		pinPairs = strings.Split(args[1], ",")
		files = args[2:]
	}

	if len(files) == 0 {
		return fmt.Errorf("at least one env file required")
	}

	pins, err := parsePinPairs(pinPairs)
	if err != nil {
		return fmt.Errorf("invalid --pins argument: %w", err)
	}

	rules := []diff.DefaultRule{diff.DefaultRuleCommonPlaceholders}
	if len(pins) > 0 {
		rules = append(rules, diff.DefaultRulePinnedValues(pins))
	}

	found := false
	for _, f := range files {
		env, err := parser.ParseFile(f)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", f, err)
		}
		issues := diff.CheckDefaults(env, rules)
		if len(issues) > 0 {
			found = true
			fmt.Fprintf(os.Stdout, "=== %s ===\n", f)
			fmt.Fprint(os.Stdout, diff.FormatDefaultIssues(issues))
		}
	}

	if !found {
		fmt.Fprintln(os.Stdout, "no default value issues found")
	}
	return nil
}

func parsePinPairs(pairs []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, p := range pairs {
		if p == "" {
			continue
		}
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected KEY=VAL, got %q", p)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}
