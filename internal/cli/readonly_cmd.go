package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunReadonly checks that keys marked as readonly have not changed between a
// reference env file and one or more target env files.
//
// Usage: envdiff readonly --keys KEY1,KEY2 <reference.env> <target1.env> [target2.env ...]
func RunReadonly(args []string, out io.Writer) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: envdiff readonly --keys KEY1,KEY2 <reference.env> <target.env> [...]")
	}

	var readonlyKeys []string
	rest := args

	if len(args) >= 2 && args[0] == "--keys" {
		for _, k := range strings.Split(args[1], ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				readonlyKeys = append(readonlyKeys, k)
			}
		}
		rest = args[2:]
	}

	if len(rest) < 2 {
		return fmt.Errorf("at least one reference and one target env file are required")
	}

	refFile := rest[0]
	targetFiles := rest[1:]

	refVars, err := parser.ParseFile(refFile)
	if err != nil {
		return fmt.Errorf("failed to parse reference file %q: %w", refFile, err)
	}

	targets := make(map[string]map[string]string, len(targetFiles))
	for _, tf := range targetFiles {
		vars, err := parser.ParseFile(tf)
		if err != nil {
			return fmt.Errorf("failed to parse target file %q: %w", tf, err)
		}
		targets[tf] = vars
	}

	issues := diff.CheckReadonly(refVars, targets, readonlyKeys)
	fmt.Fprintln(out, diff.FormatReadonlyIssues(issues))

	if len(issues) > 0 {
		os.Exit(1)
	}
	return nil
}
