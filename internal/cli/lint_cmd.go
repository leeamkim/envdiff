package cli

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// LintOptions configures the lint command.
type LintOptions struct {
	Files       []string
	NoEmpty     bool
	NoPlaceholder bool
	ForceUpper  bool
	Out         io.Writer
}

// RunLint lints one or more env files using the selected rules.
func RunLint(opts LintOptions) error {
	if len(opts.Files) == 0 {
		return fmt.Errorf("lint: at least one file is required")
	}
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	rules := buildLintRules(opts)
	if len(rules) == 0 {
		rules = []diff.LintRule{diff.LintRuleNoEmptyValues, diff.LintRuleNoDuplicatePlaceholder, diff.LintRuleKeyUppercase}
	}

	var allIssues []diff.LintIssue
	for _, f := range opts.Files {
		env, err := parser.ParseFile(f)
		if err != nil {
			return fmt.Errorf("lint: could not parse %s: %w", f, err)
		}
		result := diff.Lint(f, env, rules)
		allIssues = append(allIssues, result.Issues...)
	}

	if len(allIssues) == 0 {
		fmt.Fprintln(out, "No lint issues found.")
		return nil
	}

	sort.Slice(allIssues, func(i, j int) bool {
		if allIssues[i].File != allIssues[j].File {
			return allIssues[i].File < allIssues[j].File
		}
		return allIssues[i].Key < allIssues[j].Key
	})
	for _, issue := range allIssues {
		fmt.Fprintln(out, issue.String())
	}
	return fmt.Errorf("lint: %d issue(s) found", len(allIssues))
}

func buildLintRules(opts LintOptions) []diff.LintRule {
	var rules []diff.LintRule
	if opts.NoEmpty {
		rules = append(rules, diff.LintRuleNoEmptyValues)
	}
	if opts.NoPlaceholder {
		rules = append(rules, diff.LintRuleNoDuplicatePlaceholder)
	}
	if opts.ForceUpper {
		rules = append(rules, diff.LintRuleKeyUppercase)
	}
	return rules
}
