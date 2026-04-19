package cli

import (
	"fmt"
	"os"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunPin checks that specific key=value pairs are pinned in an env file.
// Usage: envdiff pin <file> KEY=VALUE [KEY=VALUE ...]
func RunPin(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff pin <file> KEY=VALUE [KEY=VALUE ...]")
	}

	filePath := args[0]
	pinArgs := args[1:]

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	pins, err := parsePinArgs(pinArgs)
	if err != nil {
		return err
	}

	issues := diff.CheckPins(env, pins)
	fmt.Print(diff.FormatPinIssues(issues))

	if len(issues) > 0 {
		os.Exit(1)
	}
	return nil
}

func parsePinArgs(args []string) (map[string]string, error) {
	pins := make(map[string]string, len(args))
	for _, arg := range args {
		idx := strings.IndexByte(arg, '=')
		if idx < 1 {
			return nil, fmt.Errorf("invalid pin argument %q: expected KEY=VALUE", arg)
		}
		pins[arg[:idx]] = arg[idx+1:]
	}
	return pins, nil
}
