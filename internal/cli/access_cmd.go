package cli

import (
	"fmt"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunAccess checks that keys in an env file satisfy access-level rules.
// Usage: envdiff access <file> <rule1> [rule2 ...]
// Rule format: PREFIX:LEVEL  e.g. SECRET_:secret  INTERNAL_:internal
func RunAccess(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff access <file> <PREFIX:LEVEL> [...]")
	}

	filePath := args[0]
	ruleArgs := args[1:]

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	rules, err := parseAccessRules(ruleArgs)
	if err != nil {
		return err
	}

	// Infer actual access levels from key names (heuristic).
	actual := inferActualLevels(env)

	issues := diff.CheckAccess(env, rules, actual)
	if len(issues) == 0 {
		fmt.Fprintln(os.Stdout, "all keys satisfy access policy")
		return nil
	}

	fmt.Fprint(os.Stdout, diff.FormatAccessIssues(issues))
	return fmt.Errorf("%d access violation(s) found", len(issues))
}

func parseAccessRules(args []string) ([]diff.AccessRule, error) {
	var rules []diff.AccessRule
	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rule %q: expected PREFIX:LEVEL", arg)
		}
		level, err := parseAccessLevel(parts[1])
		if err != nil {
			return nil, err
		}
		rules = append(rules, diff.AccessRule{Prefix: parts[0], Level: level})
	}
	return rules, nil
}

func parseAccessLevel(s string) (diff.AccessLevel, error) {
	switch strings.ToLower(s) {
	case "public":
		return diff.AccessPublic, nil
	case "internal":
		return diff.AccessInternal, nil
	case "secret":
		return diff.AccessSecret, nil
	default:
		return diff.AccessPublic, fmt.Errorf("unknown access level %q (use: public, internal, secret)", s)
	}
}

// inferActualLevels uses key-name patterns to guess current access level.
func inferActualLevels(env map[string]string) map[string]diff.AccessLevel {
	m := make(map[string]diff.AccessLevel, len(env))
	for k := range env {
		upper := strings.ToUpper(k)
		switch {
		case strings.Contains(upper, "SECRET") || strings.Contains(upper, "TOKEN") ||
			strings.Contains(upper, "PASSWORD") || strings.Contains(upper, "PASS") ||
			strings.Contains(upper, "PRIVATE"):
			m[k] = diff.AccessSecret
		case strings.Contains(upper, "INTERNAL") || strings.Contains(upper, "PRIV"):
			m[k] = diff.AccessInternal
		default:
			m[k] = diff.AccessPublic
		}
	}
	return m
}
