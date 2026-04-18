package cli

import (
	"fmt"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunSchema validates an env file against a required-keys schema and optional rules.
func RunSchema(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff schema <file> [--require KEY1,KEY2,...]")
	}

	filePath := args[0]
	var requiredKeys []string

	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--require=") {
			val := strings.TrimPrefix(args[i], "--require=")
			requiredKeys = strings.Split(val, ",")
		}
	}

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	rules := []diff.SchemaRule{diff.SchemaRuleNoEmptyValues}
	issues := diff.ValidateSchema(env, rules)

	if len(requiredKeys) > 0 {
		requiredRule := diff.SchemaRuleRequiredKeys(requiredKeys)
		issues = append(issues, requiredRule(env)...)
	}

	if len(issues) == 0 {
		fmt.Fprintln(os.Stdout, "schema: OK")
		return nil
	}

	for _, issue := range issues {
		fmt.Fprintln(os.Stdout, issue.String())
	}
	return fmt.Errorf("schema validation failed with %d issue(s)", len(issues))
}
