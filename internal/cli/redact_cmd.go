package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunRedact loads an env file and prints a redacted version as JSON.
// Usage: envdiff redact <file> [--patterns=A,B,...]
//
// By default, common sensitive key patterns (e.g. SECRET, PASSWORD, TOKEN)
// are redacted. Use --patterns to override with a comma-separated list.
func RunRedact(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff redact <file> [--patterns=A,B,...]")
	}

	filePath := args[0]
	patterns := diff.DefaultRedactPatterns

	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "--patterns=") {
			val := strings.TrimPrefix(arg, "--patterns=")
			if val != "" {
				patterns = strings.Split(val, ",")
			}
		}
	}

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	if len(env) == 0 {
		fmt.Fprintf(os.Stderr, "warning: %s contains no entries\n", filePath)
	}

	opts := diff.RedactOptions{Patterns: patterns}
	redacted := diff.Redact(env, opts)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(redacted); err != nil {
		return fmt.Errorf("failed to encode output: %w", err)
	}
	return nil
}
