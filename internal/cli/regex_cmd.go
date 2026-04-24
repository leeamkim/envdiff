package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/your-org/envdiff/internal/diff"
	"github.com/your-org/envdiff/internal/parser"
)

// RunRegex validates env file values against user-supplied regex rules.
//
// Usage: envdiff regex <file> KEY=PATTERN [KEY=PATTERN ...]
//
// KEY may end with '*' to match all keys sharing that prefix.
// Example: envdiff regex .env PORT=^\d+$ DB_*=^[a-z]
func RunRegex(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff regex <file> KEY=PATTERN [KEY=PATTERN ...]")
	}

	envFile := args[0]
	env, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse %q: %w", envFile, err)
	}

	rules, err := parseRegexRules(args[1:])
	if err != nil {
		return err
	}

	issues := diff.CheckRegex(env, rules)
	fmt.Print(diff.FormatRegexIssues(issues))

	if len(issues) > 0 {
		os.Exit(1)
	}
	return nil
}

// parseRegexRules converts "KEY=PATTERN" strings into RegexRule slice.
func parseRegexRules(pairs []string) ([]diff.RegexRule, error) {
	var rules []diff.RegexRule
	for _, pair := range pairs {
		idx := strings.IndexByte(pair, '=')
		if idx < 1 {
			return nil, fmt.Errorf("invalid rule %q: expected KEY=PATTERN", pair)
		}
		key := pair[:idx]
		patStr := pair[idx+1:]
		pat, err := regexp.Compile(patStr)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern %q for key %q: %w", patStr, key, err)
		}
		rules = append(rules, diff.RegexRule{Key: key, Pattern: pat})
	}
	return rules, nil
}
