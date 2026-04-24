package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunScope checks that keys in each env file conform to declared scope prefix rules.
// Usage: envdiff scope <scope:prefix1,prefix2>=<file> [...]
//
// Example:
//
//	envdiff scope prod:APP_,SVC_=.env.prod staging:APP_=.env.staging
func RunScope(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff scope <scope:prefix1,prefix2>=<file> [...]")
	}

	envs := make(map[string]map[string]string)
	rules := make([]diff.ScopeRule, 0, len(args))

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid argument %q: expected <scope:prefixes>=<file>", arg)
		}
		decl, file := parts[0], parts[1]

		declParts := strings.SplitN(decl, ":", 2)
		if len(declParts) != 2 {
			return fmt.Errorf("invalid scope declaration %q: expected <scope:prefixes>", decl)
		}
		scope := declParts[0]
		prefixes := strings.Split(declParts[1], ",")

		vars, err := parser.ParseFile(file)
		if err != nil {
			return fmt.Errorf("failed to parse %q: %w", file, err)
		}

		envs[scope] = vars
		rules = append(rules, diff.ScopeRule{Scope: scope, Prefixes: prefixes})
	}

	issues := diff.CheckScopes(envs, rules)
	fmt.Fprintln(os.Stdout, diff.FormatScopeIssues(issues))

	if len(issues) > 0 {
		return fmt.Errorf("%d scope violation(s) found", len(issues))
	}
	return nil
}
