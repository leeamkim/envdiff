package cli

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunSecrets checks an env file for secret hygiene issues.
func RunSecrets(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff secrets <file> [--weak] [--empty]")
	}

	filePath := args[0]
	checkWeak := true
	checkEmpty := true

	for _, arg := range args[1:] {
		switch arg {
		case "--weak":
			checkEmpty = false
		case "--empty":
			checkWeak = false
		}
	}

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	rules := buildSecretRules(checkWeak, checkEmpty)
	issues := diff.CheckSecrets(env, rules)

	if len(issues) == 0 {
		fmt.Fprintln(os.Stdout, "no secret issues found")
		return nil
	}

	for _, issue := range issues {
		fmt.Fprintln(os.Stdout, issue.String())
	}
	return nil
}

func buildSecretRules(weak, empty bool) []diff.SecretRule {
	var rules []diff.SecretRule
	if weak {
		rules = append(rules, diff.SecretRuleWeakValue)
	}
	if empty {
		rules = append(rules, diff.SecretRuleKeyWithoutValue)
	}
	return rules
}
