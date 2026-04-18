package cli

import (
	"flag"
	"fmt"
	"os"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

// RunTemplate compares an env file against a template file.
func RunTemplate(args []string) int {
	fs := flag.NewFlagSet("template", flag.ContinueOnError)
	strict := fs.Bool("strict", false, "fail on extra keys not in template")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}

	positional := fs.Args()
	if len(positional) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff template [--strict] <template> <env>")
		return 1
	}

	tmplMap, err := parser.ParseFile(positional[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading template: %v\n", err)
		return 1
	}

	envMap, err := parser.ParseFile(positional[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading env: %v\n", err)
		return 1
	}

	issues := diff.CompareToTemplate(tmplMap, envMap, *strict)
	fmt.Println(diff.FormatTemplateIssues(issues))

	if len(issues) > 0 {
		return 1
	}
	return 0
}
