package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunLifecycle checks an env file against lifecycle rules provided as
// key=stage pairs, e.g. OLD_*=deprecated LEGACY_TOKEN=retired.
//
// Usage: envdiff lifecycle <file> <KEY_PATTERN=stage>...
func RunLifecycle(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff lifecycle <file> <PATTERN=stage>...")
	}

	filePath := args[0]
	ruleArgs := args[1:]

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	rules, err := parseLifecycleRules(ruleArgs)
	if err != nil {
		return err
	}

	issues := diff.CheckLifecycle(env, rules)
	fmt.Println(diff.FormatLifecycleIssues(issues))

	if len(issues) > 0 {
		os.Exit(1)
	}
	return nil
}

func parseLifecycleRules(args []string) ([]diff.LifecycleRule, error) {
	var rules []diff.LifecycleRule
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid lifecycle rule %q: expected PATTERN=stage", arg)
		}
		pattern := strings.TrimSpace(parts[0])
		stageStr := strings.TrimSpace(strings.ToLower(parts[1]))

		var stage diff.LifecycleStage
		switch stageStr {
		case "active":
			stage = diff.StageActive
		case "deprecated":
			stage = diff.StageDeprecated
		case "experimental":
			stage = diff.StageExperimental
		case "retired":
			stage = diff.StageRetired
		default:
			return nil, fmt.Errorf("unknown lifecycle stage %q (use: active, deprecated, experimental, retired)", stageStr)
		}

		rules = append(rules, diff.LifecycleRule{Pattern: pattern, Stage: stage})
	}
	return rules, nil
}
