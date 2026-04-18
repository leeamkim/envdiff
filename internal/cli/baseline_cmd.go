package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunBaseline handles the baseline subcommand.
// Usage: envdiff baseline save <file> <output.json>
//        envdiff baseline check <file> <baseline.json>
func RunBaseline(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: envdiff baseline <save|check> <envfile> <jsonfile>")
	}
	subcmd, envFile, jsonFile := args[0], args[1], args[2]

	switch subcmd {
	case "save":
		return baselineSave(envFile, jsonFile)
	case "check":
		return baselineCheck(envFile, jsonFile)
	default:
		return fmt.Errorf("unknown baseline subcommand: %q", subcmd)
	}
}

func baselineSave(envFile, outFile string) error {
	env, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	b := diff.NewBaseline(env)
	data, err := diff.MarshalBaseline(b)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	if err := os.WriteFile(outFile, data, 0644); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	fmt.Printf("baseline saved to %s\n", outFile)
	return nil
}

func baselineCheck(envFile, jsonFile string) error {
	env, err := parser.ParseFile(envFile)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("read baseline error: %w", err)
	}
	baseline, err := diff.UnmarshalBaseline(data)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	d := diff.CompareToBaseline(baseline, env)
	if !d.HasDiff() {
		fmt.Println("no changes from baseline")
		return nil
	}
	keys := func(m map[string]string) []string {
		out := make([]string, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}
	for _, k := range keys(d.Added) {
		fmt.Printf("+ %s=%s\n", k, d.Added[k])
	}
	for _, k := range keys(d.Removed) {
		fmt.Printf("- %s=%s\n", k, d.Removed[k])
	}
	for k, pair := range d.Changed {
		fmt.Printf("~ %s: %q -> %q\n", k, pair[0], pair[1])
	}
	return nil
}
