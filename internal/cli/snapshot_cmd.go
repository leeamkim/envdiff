package cli

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunSnapshot handles the snapshot subcommand.
// Usage:
//   snapshot save <label> <envfile> <outfile>
//   snapshot diff <snap1> <snap2>
func RunSnapshot(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: snapshot <save|diff> [args...]")
	}
	switch args[0] {
	case "save":
		return snapshotSave(args[1:])
	case "diff":
		return snapshotDiff(args[1:])
	default:
		return fmt.Errorf("unknown snapshot command: %s", args[0])
	}
}

func snapshotSave(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: snapshot save <label> <envfile> <outfile>")
	}
	label, envPath, outPath := args[0], args[1], args[2]
	entries, err := parser.ParseFile(envPath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", envPath, err)
	}
	s := diff.NewSnapshot(label, entries)
	if err := diff.SaveSnapshot(outPath, s); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}
	fmt.Fprintf(os.Stdout, "snapshot saved: %s (%d keys)\n", outPath, len(entries))
	return nil
}

func snapshotDiff(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: snapshot diff <snap1> <snap2>")
	}
	a, err := diff.LoadSnapshot(args[0])
	if err != nil {
		return fmt.Errorf("load %s: %w", args[0], err)
	}
	b, err := diff.LoadSnapshot(args[1])
	if err != nil {
		return fmt.Errorf("load %s: %w", args[1], err)
	}
	d := diff.DiffSnapshots(a, b)
	if !d.HasChanges() {
		fmt.Fprintln(os.Stdout, "no changes between snapshots")
		return nil
	}
	for k, v := range d.Added {
		fmt.Fprintf(os.Stdout, "+ %s=%s\n", k, v)
	}
	for k, v := range d.Removed {
		fmt.Fprintf(os.Stdout, "- %s=%s\n", k, v)
	}
	for k, pair := range d.Changed {
		fmt.Fprintf(os.Stdout, "~ %s: %s -> %s\n", k, pair[0], pair[1])
	}
	return nil
}
