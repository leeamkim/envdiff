package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunClamp validates value lengths in an env file against a set of key:min:max rules.
//
// Usage: envdiff clamp <file> <KEY:min:max> [KEY:min:max ...]
func RunClamp(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff clamp <file> <KEY:min:max> [...]")
	}

	envFile := args[0]
	ruleArgs := args[1:]

	env, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	rules, err := parseClampRules(ruleArgs)
	if err != nil {
		return err
	}

	issues := diff.CheckClamps(env, rules)
	fmt.Print(diff.FormatClampIssues(issues))

	if len(issues) > 0 {
		os.Exit(1)
	}
	return nil
}

func parseClampRules(args []string) ([]diff.ClampRule, error) {
	var rules []diff.ClampRule
	for _, arg := range args {
		parts := strings.Split(arg, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid clamp rule %q: expected KEY:min:max", arg)
		}
		min, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid min in rule %q: %w", arg, err)
		}
		max, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid max in rule %q: %w", arg, err)
		}
		if min > max {
			return nil, fmt.Errorf("min > max in rule %q", arg)
		}
		rules = append(rules, diff.ClampRule{Key: parts[0], Min: min, Max: max})
	}
	return rules, nil
}
