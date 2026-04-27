package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunDependency runs the dependency check command.
// Usage: envdiff dependency <envfile> <KEY:REFKEY>...
func RunDependency(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff dependency <envfile> <KEY:REFKEY>...")
	}

	envFile := args[0]
	ruleArgs := args[1:]

	env, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", envFile, err)
	}

	rules, err := parseDependencyRules(ruleArgs)
	if err != nil {
		return err
	}

	issues := diff.CheckDependencies(env, rules)
	if len(issues) == 0 {
		fmt.Println("no dependency issues found")
		return nil
	}

	fmt.Fprint(os.Stdout, diff.FormatDependencyIssues(issues))
	return fmt.Errorf("%d dependency issue(s) found", len(issues))
}

func parseDependencyRules(args []string) ([]diff.DependencyRule, error) {
	var rules []diff.DependencyRule
	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid dependency rule %q: expected KEY:REFKEY", arg)
		}
		rules = append(rules, diff.DependencyRule{
			Key:    parts[0],
			RefKey: parts[1],
		})
	}
	return rules, nil
}
