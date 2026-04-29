package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunQuota checks that key counts per prefix do not exceed defined quotas.
// Usage: envdiff quota <file> <PREFIX:max> [<PREFIX:max> ...]
func RunQuota(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff quota <file> <PREFIX:max> [<PREFIX:max> ...]")
	}

	envFile := args[0]
	ruleArgs := args[1:]

	env, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", envFile, err)
	}

	rules, err := parseQuotaRules(ruleArgs)
	if err != nil {
		return err
	}

	issues := diff.CheckQuotas(env, rules)
	fmt.Print(diff.FormatQuotaIssues(issues))

	if len(issues) > 0 {
		os.Exit(1)
	}
	return nil
}

func parseQuotaRules(args []string) ([]diff.QuotaRule, error) {
	var rules []diff.QuotaRule
	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid quota rule %q: expected PREFIX:max", arg)
		}
		max, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid max value in rule %q: %w", arg, err)
		}
		if max < 0 {
			return nil, fmt.Errorf("max value must be non-negative in rule %q", arg)
		}
		rules = append(rules, diff.QuotaRule{Prefix: parts[0], MaxKeys: max})
	}
	return rules, nil
}
