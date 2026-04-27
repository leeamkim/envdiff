package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunImmutable checks that keys declared as immutable have not changed across envs.
// Usage: envdiff immutable --keys KEY1,KEY2 env1=file1.env env2=file2.env ...
func RunImmutable(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff immutable --keys KEY1,KEY2 name=file ...")
	}

	var immutableKeys []string
	var envArgs []string

	for i := 0; i < len(args); i++ {
		if args[i] == "--keys" && i+1 < len(args) {
			immutableKeys = strings.Split(args[i+1], ",")
			i++
		} else {
			envArgs = append(envArgs, args[i])
		}
	}

	if len(immutableKeys) == 0 {
		return fmt.Errorf("--keys is required and must not be empty")
	}
	if len(envArgs) < 2 {
		return fmt.Errorf("at least two name=file arguments are required")
	}

	envs := make(map[string]map[string]string, len(envArgs))
	for _, arg := range envArgs {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid argument %q: expected name=file", arg)
		}
		name, file := parts[0], parts[1]
		kvs, err := parser.ParseFile(file)
		if err != nil {
			return fmt.Errorf("parsing %q: %w", file, err)
		}
		envs[name] = kvs
	}

	issues := diff.CheckImmutable(envs, immutableKeys)
	fmt.Fprint(os.Stdout, diff.FormatImmutableIssues(issues))

	if len(issues) > 0 {
		return fmt.Errorf("immutable violations detected")
	}
	return nil
}
