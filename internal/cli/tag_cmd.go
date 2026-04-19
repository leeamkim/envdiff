package cli

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunTag runs the tag command, tagging keys in an env file by rule.
func RunTag(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff tag <file> [--secrets] [--empty] [--placeholder]")
	}

	filePath := args[0]
	flags := map[string]bool{}
	for _, a := range args[1:] {
		flags[a] = true
	}

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	var rules []diff.TagRule
	if flags["--secrets"] {
		rules = append(rules, diff.TagRuleSecrets)
	}
	if flags["--empty"] {
		rules = append(rules, diff.TagRuleEmpty)
	}
	if flags["--placeholder"] {
		rules = append(rules, diff.TagRulePlaceholder)
	}
	if len(rules) == 0 {
		rules = []diff.TagRule{diff.TagRuleSecrets, diff.TagRuleEmpty, diff.TagRulePlaceholder}
	}

	entries := diff.TagEnv(env, rules...)
	fmt.Fprint(os.Stdout, diff.FormatTagEntries(entries))
	return nil
}
