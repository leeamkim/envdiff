package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunCrossRef checks cross-references within a single .env file.
// Usage: envdiff crossref <file> [--rule SRC:TARGET] ...
func RunCrossRef(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff crossref <file> [--rule SRC:TARGET]...")
	}

	envFile := args[0]
	env, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", envFile, err)
	}

	rules, err := parseCrossRefRules(args[1:])
	if err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	envName := envFile
	issues := diff.CheckCrossRefs(env, rules, envName)

	if len(issues) == 0 {
		fmt.Fprintln(os.Stdout, "no cross-reference issues found")
		return nil
	}

	fmt.Fprintln(os.Stdout, diff.FormatCrossRefIssues(issues))
	return fmt.Errorf("%d cross-reference issue(s) found", len(issues))
}

// parseCrossRefRules parses --rule SRC:TARGET flags from args.
func parseCrossRefRules(args []string) ([]diff.CrossRefRule, error) {
	var rules []diff.CrossRefRule
	for i := 0; i < len(args); i++ {
		if args[i] == "--rule" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--rule requires an argument")
			}
			i++
			parts := strings.SplitN(args[i], ":", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				return nil, fmt.Errorf("rule %q must be in SRC:TARGET format", args[i])
			}
			rules = append(rules, diff.CrossRefRule{
				SourceKey: parts[0],
				TargetKey: parts[1],
			})
		}
	}
	return rules, nil
}
