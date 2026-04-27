package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// RunExpire checks env file keys against expiry rules.
// Usage: envdiff expire <file> <KEY_PATTERN=YYYY-MM-DD,...> [--warn-days=N]
func RunExpire(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff expire <file> <KEY=YYYY-MM-DD,...> [--warn-days=N]")
	}

	filePath := args[0]
	rulesArg := args[1]
	warnDays := 30

	for _, arg := range args[2:] {
		if strings.HasPrefix(arg, "--warn-days=") {
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "--warn-days="))
			if err != nil {
				return fmt.Errorf("invalid --warn-days value: %w", err)
			}
			warnDays = n
		}
	}

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	rules, err := parseExpiryRules(rulesArg)
	if err != nil {
		return fmt.Errorf("failed to parse expiry rules: %w", err)
	}

	issues := diff.CheckExpiry(env, rules, warnDays, time.Now().UTC())
	fmt.Fprintln(os.Stdout, diff.FormatExpireIssues(issues))
	return nil
}

func parseExpiryRules(raw string) ([]diff.ExpiryRule, error) {
	var rules []diff.ExpiryRule
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx := strings.LastIndex(part, "=")
		if idx < 0 {
			return nil, fmt.Errorf("invalid rule %q: expected KEY=YYYY-MM-DD", part)
		}
		pattern := part[:idx]
		dateStr := part[idx+1:]
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date %q in rule %q: %w", dateStr, part, err)
		}
		rules = append(rules, diff.ExpiryRule{Pattern: pattern, ExpiresAt: t})
	}
	return rules, nil
}
