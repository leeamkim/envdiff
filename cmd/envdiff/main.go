package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// version information set by the build system via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)
