package cli

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunAudit compares two env files and prints an audit log of changes.
func RunAudit(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff audit <before.env> <after.env>")
	}

	before, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("reading before file: %w", err)
	}

	after, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("reading after file: %w", err)
	}

	log := diff.GenerateAuditLog(before, after)

	if !log.HasChanges() {
		fmt.Fprintln(os.Stdout, "No changes detected.")
		return nil
	}

	for _, entry := range log {
		fmt.Fprintln(os.Stdout, entry.String())
	}
	return nil
}
