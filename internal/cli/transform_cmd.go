package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunTransform applies transformations to an env file and prints the result.
// Usage: envdiff transform <file> --op=<trim|upper|lower> [--keys=A,B,...]
func RunTransform(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff transform <file> [--op=trim|upper|lower] [--keys=K1,K2]")
	}
	file := args[0]
	op := "trim"
	var keys []string
	for _, a := range args[1:] {
		if strings.HasPrefix(a, "--op=") {
			op = strings.TrimPrefix(a, "--op=")
		} else if strings.HasPrefix(a, "--keys=") {
			raw := strings.TrimPrefix(a, "--keys=")
			for _, k := range strings.Split(raw, ",") {
				if k != "" {
					keys = append(keys, k)
				}
			}
		}
	}
	env, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parse %s: %w", file, err)
	}
	var fn diff.TransformFunc
	switch op {
	case "trim":
		fn = diff.TransformTrimSpace()
	case "upper":
		fn = diff.TransformUppercase()
	case "lower":
		fn = diff.TransformLowercase()
	default:
		return fmt.Errorf("unknown op %q: choose trim, upper, or lower", op)
	}
	res := diff.Transform(env, []diff.TransformOption{{Keys: keys, Fn: fn}})
	if len(res.Changed) == 0 {
		fmt.Fprintln(os.Stdout, "No changes.")
		return nil
	}
	fmt.Fprintf(os.Stdout, "Changed %d key(s):\n", len(res.Changed))
	for _, k := range res.Changed {
		fmt.Fprintf(os.Stdout, "  %s=%s\n", k, res.Out[k])
	}
	return nil
}
