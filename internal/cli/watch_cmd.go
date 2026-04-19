package cli

import (
	"fmt"
	"io"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// WatchOptions configures the watch command.
type WatchOptions struct {
	Files    []string
	Interval time.Duration
	MaxTicks int // 0 = unlimited; used for testing
}

// RunWatch polls the given env files and prints a diff whenever a change is detected.
func RunWatch(opts WatchOptions, out io.Writer) error {
	if len(opts.Files) < 2 {
		return fmt.Errorf("watch requires at least two files")
	}
	if opts.Interval <= 0 {
		opts.Interval = 2 * time.Second
	}

	hashes := map[string]string{}
	// Seed initial hashes.
	for _, f := range opts.Files {
		h, err := diff.HashFile(f)
		if err != nil {
			return err
		}
		hashes[f] = h
	}

	tick := 0
	for {
		time.Sleep(opts.Interval)
		events, updated, err := diff.WatchOnce(opts.Files, hashes)
		if err != nil {
			return err
		}
		hashes = updated
		if len(events) > 0 {
			for _, e := range events {
				fmt.Fprintln(out, e.String())
			}
			a, err := parser.ParseFile(opts.Files[0])
			if err != nil {
				return err
			}
			b, err := parser.ParseFile(opts.Files[1])
			if err != nil {
				return err
			}
			result := diff.Compare(a, b)
			diff.PrintReport(result, out)
		}
		tick++
		if opts.MaxTicks > 0 && tick >= opts.MaxTicks {
			break
		}
	}
	return nil
}
